package handlers

import (
	"net/http"

	"github.com/CRoasSanhez/yofio-test/internal/platform/web"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// API REST interface for api
func API(db *gorm.DB) http.Handler {
	srv, err := web.NewServer(db)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error creating a new server")
	}

	srv.HandleService("POST", "/api/register", Register)
	srv.HandleService("POST", "/api/login", Login)
	srv.HandleService("POST", "/api/payments", MembershipPayment)
	srv.HandleService("GET", "/api/payments", MembershipConsult)

	return srv
}
