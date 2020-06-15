package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/CRoasSanhez/yofio-test/internal/platform/auth"

	"github.com/CRoasSanhez/yofio-test/internal/platform/database/schema"
	"github.com/CRoasSanhez/yofio-test/internal/utils"
	"github.com/sirupsen/logrus"

	"github.com/CRoasSanhez/yofio-test/internal/platform/web"
	"github.com/julienschmidt/httprouter"
)

// Register ...
func Register(server *web.Server) web.ServerHandler {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {

		var reqUser = &struct {
			FullName string `json:"name"`
			Phone    string `json:"phone"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}

		if err := web.DecodeJSONRequest(r, &reqUser); err != nil {
			logrus.WithFields(logrus.Fields{
				"decodeUserRequestErr": err,
			}).Error(err.Error())

			return web.ErrorInvalidData.WithStatus(http.StatusBadRequest) // 400
		}

		var user = &schema.User{}
		if err := utils.CopyStruct(reqUser, user); err != nil {
			logrus.WithFields(logrus.Fields{
				"decodeUserStructtErr": err,
			}).Error(err.Error())

			return web.ErrorInvalidData.WithStatus(http.StatusBadRequest) // 400
		}

		if !utils.IsValidEmail(user.Email) ||
			!utils.IsValidPhone(user.Phone) ||
			!user.IsValidFullName() ||
			!user.IsValidPassword() {
			logrus.WithFields(logrus.Fields{
				"email": user.Email,
				"phone": user.Phone,
				"name":  user.FullName,
			}).Info()
			return web.ErrorInvalidData.WithStatus(http.StatusBadRequest) // 400
		}

		// VALIDATE UNIQUE USER PHONE IN DB
		var userDB = &schema.User{}
		server.DataBase.Where("Phone=?", user.Phone).First(userDB)
		if userDB.Email != "" {
			logrus.WithFields(logrus.Fields{
				"Phone already in use": user.Phone,
			}).Info()
			return web.PhoneRepeated.WithStatus(http.StatusPreconditionFailed) // 412
		}

		user.Pwd = utils.HashPassword(reqUser.Password)

		server.DataBase.Create(user)

		web.ResponseJSON(w, http.StatusOK, &struct{}{})
		return nil
	}
}

// Login ...
func Login(server *web.Server) web.ServerHandler {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {

		var userReq = &struct {
			Phone    string `json:"phone"`
			Password string `json:"password"`
		}{}

		if err := web.DecodeJSONRequest(r, &userReq); err != nil {
			logrus.WithFields(logrus.Fields{
				"decodeUserRequestErr": err,
			}).Error(err.Error())

			return web.ErrorInvalidData.WithStatus(http.StatusBadRequest)
		}

		// Validate email regex
		if !utils.IsValidPhone(userReq.Phone) {
			logrus.WithFields(logrus.Fields{
				"invalid phone": userReq.Phone,
			}).Info()
			return web.PhoneRepeated.WithStatus(http.StatusPreconditionFailed) // 412
		}

		// Retreive user from DB
		userDB := &schema.User{}
		server.DataBase.Where("Phone=?", userReq.Phone).First(userDB)

		// Validate blocked user
		if userDB.IsBlocked {
			logrus.WithFields(logrus.Fields{
				"User blocked": userDB.ID,
			}).Info()
			return web.UserBlocked.WithStatus(http.StatusForbidden) //403
		}

		// Validate login failures
		if userDB.LoginFailures >= 5 {
			logrus.WithFields(logrus.Fields{
				"Max attemps reached": userDB.ID,
			}).Info()
			return web.MaxLoginAttemptsReached.WithStatus(http.StatusTooManyRequests) //429
		}

		// Verify password
		if !utils.CheckPassword(userDB.Pwd, userReq.Password) {
			logrus.WithFields(logrus.Fields{
				"invalid Password": userReq.Password,
			}).Info()

			// Update Attempts in DB
			go userDB.AddLoginAttempts(server.DataBase)
			return web.PasswordInvalid.WithStatus(http.StatusUnauthorized) //401
		}

		// Generate Token for loged in user
		token, err := auth.SignToken(strconv.Itoa(userDB.ID))
		if err != nil {

			logrus.WithFields(logrus.Fields{
				"error":  err.Error(),
				"userID": userDB.ID,
			}).Info("Error Generating Token")

			return web.ErrorGenerateJWT.WithStatus(http.StatusInternalServerError)
		}

		web.ResponseJSON(w, http.StatusOK, &struct {
			Token string `json:"token"`
		}{
			Token: token,
		})

		return nil

	}
}

// MembershipPayment ...
func MembershipPayment(server *web.Server) web.ServerHandler {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {

		user, err := auth.GetCurrentUser(server.DataBase, r.Header.Get("Authorization"))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"user": user.Email,
			}).Error(err.Error())
			return web.ErrorLoginUser.WithStatus(http.StatusUnauthorized) // 401
		}

		// Validate blocked user
		if user.IsBlocked {
			logrus.WithFields(logrus.Fields{
				"user blocked": user.Email,
			}).Error(err.Error())
			return web.UserBlocked.WithStatus(http.StatusForbidden) // 403
		}

		var payReq = &struct {
			Amount     int    `json:"amount"`
			CardNumber string `json:"card_number"`
			ExpDate    string `json:"exp_date"`
			Owner      string `json:"owner"`
			CVV        string `json:"cvv"`
		}{}

		if err := web.DecodeJSONRequest(r, &payReq); err != nil {
			logrus.WithFields(logrus.Fields{
				"decodePaymentRequestErr": err,
			}).Error(err.Error())

			return web.ErrorInvalidData.WithStatus(http.StatusBadRequest)
		}

		// Find Current User membership
		var membership = &schema.Membership{}
		server.DataBase.Where("ID=?", user.ID).First(membership)

		// Create membership if does not exist
		if membership.CardNumber == "" {
			membership = schema.NewMembership(user.ID, 100000, payReq.CardNumber, payReq.Owner)
			server.DataBase.Create(membership)
		}

		membership.Attempts++

		var last4digits = string(payReq.CardNumber[len(payReq.CardNumber)-4])

		// Validate incoming data
		if len(payReq.CVV) != 3 ||
			membership.CardNumber != payReq.CardNumber ||
			!utils.IsValidAmunt(payReq.Amount) ||
			!utils.IsValidCard(payReq.CardNumber) ||
			!utils.IsValidExpirationDate(payReq.ExpDate) {
			logrus.WithFields(logrus.Fields{
				"error":      "Invalid data",
				"cvv length": len(payReq.CVV),
				"amount":     payReq.Amount,
				"card":       fmt.Sprintf("**** %s", last4digits),
				"expDate":    utils.IsValidExpirationDate(payReq.ExpDate),
			}).Error()

			return web.ErrorInvalidData.WithStatus(membership.SaveAttempt(server.DataBase, payReq.Amount)) // 400
		}

		var lastNumber, errAtoi = strconv.Atoi(string(last4digits[len(last4digits)-1]))
		if errAtoi != nil {
			logrus.WithFields(logrus.Fields{
				"strconvErr": payReq.CardNumber,
			}).Error(errAtoi.Error())
			return web.ErrorUpdatingMembership.WithStatus(http.StatusInternalServerError)
		}

		if utils.IsOdd(lastNumber) {
			membership.SaveAttempt(server.DataBase, payReq.Amount)
			return web.UnsuccessPayment.WithStatus(http.StatusUnprocessableEntity) // 422
		}

		// Save success payment
		p := &schema.Payment{
			MembershipID: membership.ID,
			UserID:       membership.UserID,
			Status:       "success",
			Amount:       payReq.Amount,
			InsertedAt:   time.Now(),
		}
		go p.SavePayment(server.DataBase)

		web.ResponseJSON(w, http.StatusOK, &struct{}{})

		return nil
	}
}

// MembershipConsult ...
func MembershipConsult(server *web.Server) web.ServerHandler {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {

		user, err := auth.GetCurrentUser(server.DataBase, r.Header.Get("Authorization"))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"user": user.Email,
			}).Error(err.Error())
			return web.ErrorLoginUser.WithStatus(http.StatusUnauthorized) // 401
		}

		// Validate blocked user
		if user.IsBlocked {
			logrus.WithFields(logrus.Fields{
				"user blocked": user.Email,
			}).Error(err.Error())
			return web.UserBlocked.WithStatus(http.StatusForbidden) // 403
		}

		// Return current user payments
		payments := &[]schema.Payment{}
		server.DataBase.Where("UserID=? AND Status=?", user.ID, "success").Find(payments)

		web.ResponseJSON(w, http.StatusOK, payments)
		return nil

	}
}
