package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// Envs it is used to store enviroment vars
type Envs struct {
	Env               string `required:"true" split_words:"true"`
	AppHost           string `required:"true" split_words:"true" default:"0.0.0.0"`
	AppPort           string `required:"true" split_words:"true"`
	JwtSecretKey      string `required:"true" split_words:"true" default:"sup3rs3cret"`
	MysqlDBPassword   string `required:"true" split_words:"true"`
	MysqlDatabaseName string `required:"true" split_words:"true"`
}

// SetUpEnvs ...
func SetUpEnvs() *Envs {
	var envs Envs

	configError := envconfig.Process(os.Getenv("APPNAME"), &envs)
	if configError != nil {
		logrus.WithFields(logrus.Fields{
			"error": configError.Error(),
		}).Fatal("Error with the config")
	}

	return &envs
}
