package auth

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/CRoasSanhez/yofio-test/internal/platform/database/schema"

	jwtv4 "github.com/dgrijalva/jwt-go/v4"
	uuidv4 "github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Claims for JWT
type Claims struct {
	UserID string `db:"bot" json:"bot"`
	JTI    string `db:"jti" json:"jti"`
	jwtv4.StandardClaims
}

var jwtKey = []byte(os.Getenv("YOFIOTEST_JWTSECRET_KEY"))

const jwtIssuer = "yofio-test"
const audience = "yofio"

// SignToken always does something
func SignToken(userID string) (string, error) {
	claims := make(jwtv4.MapClaims)
	now := time.Now().Unix()

	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // expiration for 72hrs
	claims["aud"] = audience
	claims["uid"] = userID
	claims["iat"] = now
	claims["iss"] = jwtIssuer
	claims["jti"] = uuidv4.New() // RFC4122
	claims["nbf"] = now
	claims["typ"] = "access"

	var token = jwtv4.NewWithClaims(jwtv4.SigningMethodHS512, claims)

	return token.SignedString(jwtKey)

}

// GetCurrentUser ...
func GetCurrentUser(db *gorm.DB, header string) (*schema.User, error) {
	var claims = &Claims{}
	var user = &schema.User{}
	err := ValidateToken(header, claims)
	if err != nil {
		return user, err
	}

	// Get User from DB
	db.Where("ID=?", claims.UserID).First(user)

	if user.Email == "" {
		return user, errors.New("User not found")
	}

	return user, nil
}

// ValidateToken validates Authorization header for token
func ValidateToken(header string, claims *Claims) error {

	tokenStr, err := getTokenFromHeader(header)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err.Error(),
		}).Error("Token not found")
		return err
	}

	tkn, err := jwtv4.ParseWithClaims(tokenStr, claims, func(token *jwtv4.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwtv4.ErrSignatureInvalid {
			logrus.WithFields(logrus.Fields{
				"Error": err.Error(),
			}).Error("Error Invalid signature")
			return err
		}

		if len(claims.Audience) > 0 && strings.Contains(err.Error(), "audience value was expected") {
			tkn.Valid = true
		} else {
			logrus.WithFields(logrus.Fields{
				"Error": err.Error(),
			}).Error("Error Audience not provided bad request")
			return err
		}
	}

	if !tkn.Valid {
		logrus.WithFields(logrus.Fields{
			"error": "Invalid token",
		}).Error("Error Invalid token")
		return err
	}

	return nil
}

// RenewToken renews current token for user
func RenewToken(header string, claims *Claims) (string, error) {

	err := ValidateToken(header, claims)
	if err != nil {
		return "", errors.New("Invalid Token")
	}

	// Renew token when current is within 60 secs of creation
	if claims.ExpiresAt != nil && (time.Unix(claims.ExpiresAt.Unix(), 0).Sub(time.Now()) > 60*time.Second) {
		return "", errors.New("Not suitable")
	}

	userID := "sample_userID"
	tknStr, err := SignToken(userID)
	return tknStr, err
}

// getTokenFromHeader returns token from Authorization header
func getTokenFromHeader(header string) (string, error) {

	var inToken = strings.Split(header, " ")

	if header == "" || !strings.HasPrefix(inToken[0], "Bearer") || inToken[1] == "" {
		return "", errors.New("No token provided")
	}
	return inToken[1], nil
}

//ContainsAudience checks whether a given string is included in the Audience
func (c *Claims) ContainsAudience(v string) bool {
	for _, a := range c.Audience {
		if a == v {
			return true
		}
	}
	return false
}
