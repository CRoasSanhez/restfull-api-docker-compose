package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/CRoasSanhez/yofio-test/cmd/api/handlers"
	"github.com/CRoasSanhez/yofio-test/internal/config"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := run(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(err)

		os.Exit(1)
	}
}

func run() error {
	appConfig := config.SetUpEnvs()

	logrus.SetFormatter(&logrus.JSONFormatter{})

	database := initDatabase(appConfig)

	initSDKs(appConfig)

	initServer(database, appConfig)

	return nil
}

func initServer(db *sql.DB, appConfig *config.Envs) {
	httpServer := http.Server{
		Addr:    fmt.Sprintf("%s:%s", appConfig.AppHost, appConfig.AppPort),
		Handler: handlers.API(db),
	}

	logrus.WithFields(logrus.Fields{
		"host": appConfig.AppHost,
		"port": appConfig.AppPort,
	}).Info("Listening for requests")

	if err := httpServer.ListenAndServe(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error starting the server")
	}
}

func initDatabase(appConfig *config.Envs) *sql.DB {

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", appConfig.DBUser, appConfig.DBPassword, appConfig.DBName))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error starting the server")
		panic(err)
	}

	return db
}
