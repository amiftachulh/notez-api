package main

import (
	"notez-api/config"
	"notez-api/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.Setup()
	app := fiber.New()
	route.Setup(app)
	app.Listen(":3000")
}
