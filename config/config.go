package config

import (
	"log"

	"github.com/izsal/go-refresh-token/database"
	"github.com/joho/godotenv"
)

func LoadConfig() {
	err := godotenv.Load()
	database.Connect()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
