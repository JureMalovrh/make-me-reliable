package config

import (
	"os"
)

// Config is a config definition
type Config struct {
	DatabaseName string
	DatabaseURL  string
	DatabaseUser string
	DatabasePass string

	APIAuthHeader string
	APIUrl        string

	ServerPort string
}

// ParseFromEnv returns Config object parsed from env variables
// note: done by hand as it is small env otherwise use lib
func ParseFromEnv() Config {
	databaseName := os.Getenv("DATABASE_NAME")
	databaseUser := os.Getenv("DATABASE_USERNAME")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databaseUrl := os.Getenv("DATABASE_URL")

	apiAuthHeader := os.Getenv("API_AUTH")
	apiUrl := os.Getenv("API_URL")
	serverPort := os.Getenv("SERVER_PORT")

	c := Config{
		DatabaseName:  databaseName,
		DatabaseUser:  databaseUser,
		DatabasePass:  databasePassword,
		DatabaseURL:   databaseUrl,
		APIAuthHeader: apiAuthHeader,
		APIUrl:        apiUrl,
		ServerPort:    serverPort,
	}

	return c
}
