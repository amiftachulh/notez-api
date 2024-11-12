package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	DatabaseURL    string
	AllowedOrigins map[string]struct{}
)

func Setup() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	DatabaseURL = os.Getenv("DATABASE_URL")
	AllowedOrigins = make(map[string]struct{})
	for _, origin := range strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",") {
		AllowedOrigins[origin] = struct{}{}
	}
}
