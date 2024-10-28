package config

import (
	"log"
	"notez-api/db"

	"github.com/joho/godotenv"
)

func Setup() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	db.Setup()
}
