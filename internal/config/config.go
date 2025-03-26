package config

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func Get(key string) (string, error) {
	value := os.Getenv(key)

	if value == "" {
		return "", errors.New("no value for provided key")
	}

	return value, nil
}
