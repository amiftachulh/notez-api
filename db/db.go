package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func Setup() {
	var err error
	DB, err = sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln("Failed to create database connection:", err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatalln("Failed to ping the database:", err)
	}
}
