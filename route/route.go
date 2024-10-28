package route

import (
	"notez-api/handler"
	"notez-api/middleware"
	"notez-api/schema"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	v1 := app.Group("/v1")

	auth := v1.Group("/auth")
	auth.Post("/register", middleware.Validate(new(schema.Register)), handler.Register)
	auth.Post("/login", middleware.Validate(new(schema.Login)), handler.Login)
	auth.Post("/logout", handler.Logout)
	auth.Get("/check", handler.CheckAuth)
}
