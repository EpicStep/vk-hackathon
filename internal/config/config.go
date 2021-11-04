package config

import (
	"fmt"
	"os"
)

// Config struct.
type Config struct {
	Port      string
	MySQLURL     string
}

// New config.
func New() (*Config, error) {
	port, err := getEnv("PORT")
	if err != nil {
		return nil, err
	}

	mySQLURL, err := getEnv("DATABASE_URL")
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:      port,
		MySQLURL:  mySQLURL,
	}, nil
}

func getEnv(key string) (string, error) {
	value, isFounded := os.LookupEnv(key)
	if isFounded {
		return value, nil
	}

	return "", fmt.Errorf("env variable %s not presented", key)
}