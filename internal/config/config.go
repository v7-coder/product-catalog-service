package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppPort string
	DBUrl   string
}

func Load() (*Config, error) {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	switch "" {
	case dbHost:
		return nil, fmt.Errorf("DB_HOST is required")
	case dbPort:
		return nil, fmt.Errorf("DB_PORT is required")
	case dbUser:
		return nil, fmt.Errorf("DB_USER is required")
	case dbName:
		return nil, fmt.Errorf("DB_NAME is required")
	}

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	return &Config{
		AppPort: appPort,
		DBUrl:   dbUrl,
	}, nil
}
