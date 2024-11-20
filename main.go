package main

import (
	"log"

	"github.com/amiftachulh/notez-api/config"
	"github.com/amiftachulh/notez-api/db"
	"github.com/amiftachulh/notez-api/handler"
	"github.com/amiftachulh/notez-api/route"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.Setup()
	db.Setup()

	app := fiber.New(fiber.Config{
		ErrorHandler: handler.ErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			_, exists := config.AllowedOrigins[origin]
			return exists
		},
		AllowCredentials: true,
	}))

	route.Setup(app)
	log.Fatal(app.Listen("127.0.0.1:3000"))
}
