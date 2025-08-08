package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	PORT string

	DB_HOST      string
	DB_USER      string
	DB_PASSWORD  string
	DB_PORT      string
	DB_NAME      string
	REDIS_SERVER string

	JWT_ACCESS_SECRET  string
	JWT_REFRESH_SECRET string

	FROM_EMAIL    string
	SMTP_HOST     string
	SMTP_PORT     string
	SMTP_USERNAME string
	SMTP_PASSWORD string

	RABBITMQ_SERVER string
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	} else {
		fmt.Println("Loaded .env file")
	}
}

func GetEnv() Env {
	return Env{
		PORT:               os.Getenv("PORT"),
		DB_HOST:            os.Getenv("DB_HOST"),
		DB_USER:            os.Getenv("DB_USER"),
		DB_PASSWORD:        os.Getenv("DB_PASSWORD"),
		DB_PORT:            os.Getenv("DB_PORT"),
		DB_NAME:            os.Getenv("DB_NAME"),
		REDIS_SERVER:       os.Getenv("REDIS_SERVER"),
		JWT_ACCESS_SECRET:  os.Getenv("JWT_ACCESS_SECRET"),
		JWT_REFRESH_SECRET: os.Getenv("JWT_REFRESH_SECRET"),
		FROM_EMAIL:         os.Getenv("FROM_EMAIL"),
		SMTP_HOST:          os.Getenv("SMTP_HOST"),
		SMTP_PORT:          os.Getenv("SMTP_PORT"),
		SMTP_USERNAME:      os.Getenv("SMTP_USERNAME"),
		SMTP_PASSWORD:      os.Getenv("SMTP_PASSWORD"),
		RABBITMQ_SERVER:    os.Getenv("RABBITMQ_SERVER"),
	}
}
