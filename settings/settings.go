package settings

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var ENVIRONMENT string
var JWT_SECRET_KEY string
var DATABASE_DB_NAME string
var DATABASE_DB_PASSWORD string
var DATABASE_DB_USER string
var DATABASE_DB_HOST string

func init() {
	if strings.HasSuffix(os.Args[0], ".test") {

		ENVIRONMENT = "local"
		JWT_SECRET_KEY = "default-secret-key-for-testing"
		DATABASE_DB_NAME = "test"
		DATABASE_DB_PASSWORD = "password"
		DATABASE_DB_USER = "postgres"
		DATABASE_DB_HOST = "localhost"

	} else {

		ENVIRONMENT = GetEnvironmentVariableOrPanic("ENVIRONMENT")
		JWT_SECRET_KEY = GetEnvironmentVariableOrPanic("JWT_SECRET_KEY")
		DATABASE_DB_NAME = GetEnvironmentVariableOrPanic("DATABASE_DB_NAME")
		DATABASE_DB_PASSWORD = GetEnvironmentVariableOrPanic("DATABASE_DB_PASSWORD")
		DATABASE_DB_USER = GetEnvironmentVariableOrPanic("DATABASE_DB_USER")
		DATABASE_DB_HOST = GetEnvironmentVariableOrPanic("DATABASE_DB_HOST")

	}
}

func GetEnvironmentVariableOrPanic(key string) string {
	viper.SetConfigFile("environment.env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("FAILED LOADING: %s, wrong type", key)
	} else if value == "" {
		log.Fatalf("FAILED LOADING: %s, missing/empty value", key)
	}

	log.Println("SETTING LOADED -", key)
	return value
}

func GetEnvironmentVariableWithDefault(key, defaultValue string) string {
	viper.SetConfigFile("environment.env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	value, ok := viper.Get(key).(string)
	if !ok || value == "" {
		value = defaultValue
	}

	log.Println("SETTING LOADED -", key)
	return value
}
