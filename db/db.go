package db

import (
	"database/sql"
	"log"

	"github.com/amiftachulh/notez-api/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func Setup() {
	var err error
	DB, err = sql.Open("pgx", config.DatabaseURL)
	if err != nil {
		log.Fatalln("Failed to create database connection:", err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatalln("Failed to ping the database:", err)
	}
}
