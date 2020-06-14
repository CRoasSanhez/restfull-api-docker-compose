package handlers

import (
	"database/sql"
	"net/http"

	"github.com/CRoasSanhez/yofio-test/internal/platform/web"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

// API REST interface for api
func API(db *sql.DB) http.Handler {
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
