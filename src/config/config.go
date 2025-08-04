package config

import (
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
)

var NRc *newrelic.Application

var Logger *logrus.Logger
var LoggerEntry *logrus.Entry

func init() {
	pwd, _ := os.Getwd()
	viper.SetConfigFile(fmt.Sprintf("%s/.env", pwd))
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error saat membaca file .envUser: %v", err)
	}

	name := viper.GetString("LOG_NAME")
	f, err := os.OpenFile(fmt.Sprintf("%s.log", name), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	Logger = logrus.New()
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetOutput(f)
	LoggerEntry = logrus.NewEntry(Logger)

	NRc, _ = newrelic.NewApplication(
		newrelic.ConfigAppName(viper.GetString("NEW_RELIC_APP_NAME")),
		newrelic.ConfigLicense(viper.GetString("NEW_RELIC_LICENSE")),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
}
