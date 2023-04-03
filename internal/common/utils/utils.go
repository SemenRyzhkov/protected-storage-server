package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// GetServerAddress получение переменной SERVER_ADDRESS
func GetServerAddress() string {
	return os.Getenv("SERVER_ADDRESS")
}

// GetKey получение переменной SECRET_KEY
func GetKey() string {
	return os.Getenv("SECRET_KEY")
}

// GetDBAddress получение переменной DATABASE_DSN
func GetDBAddress() string {
	return os.Getenv("DATABASE_DSN")
}

// LoadEnvironments загрузка конфигов из .env файла
func LoadEnvironments(envFilePath string) {
	err := godotenv.Load(envFilePath)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
