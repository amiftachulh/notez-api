package main

import (
	"log"

	"github.com/amiftachulh/notez-api/config"
	"github.com/amiftachulh/notez-api/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.Setup()
	app := fiber.New()
	route.Setup(app)
	log.Fatal(app.Listen(":3000"))
}
