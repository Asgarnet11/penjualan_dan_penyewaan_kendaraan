// File: internal/config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBSource      string
	JWTSecretKey  string
	AppPort       string
	CloudinaryURL string
	FrontendURL   string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	return Config{
		DBSource:      os.Getenv("DB_SOURCE"),
		JWTSecretKey:  os.Getenv("JWT_SECRET_KEY"),
		AppPort:       os.Getenv("APP_PORT"),
		CloudinaryURL: os.Getenv("CLOUDINARY_URL"),
		FrontendURL:   os.Getenv("FRONTEND_URL"),
	}
}
