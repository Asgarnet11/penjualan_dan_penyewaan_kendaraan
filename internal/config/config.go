package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config menampung semua konfigurasi aplikasi
type Config struct {
	DBSource     string
	JWTSecretKey string
	AppPort      string
}

// LoadConfig membaca konfigurasi dari file .env
func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	return Config{
		DBSource:     os.Getenv("DB_SOURCE"),
		JWTSecretKey: os.Getenv("JWT_SECRET_KEY"),
		AppPort:      os.Getenv("APP_PORT"),
	}
}
