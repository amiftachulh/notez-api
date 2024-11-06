package main

import (
	"log"

	"github.com/amiftachulh/notez-api/config"
	"github.com/amiftachulh/notez-api/handler"
	"github.com/amiftachulh/notez-api/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.Setup()
	app := fiber.New(fiber.Config{
		ErrorHandler: handler.ErrorHandler,
	})
	route.Setup(app)
	log.Fatal(app.Listen(":3000"))
}
