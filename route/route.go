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

	notes := v1.Group("/notes").Use(middleware.Authenticate)

	// notes := protected.Group("/notes")
	notes.Post("/", middleware.Validate(new(schema.Note)), handler.CreateNote)
	notes.Get("/", handler.GetNotes)
	notes.Get("/:id", handler.GetNoteByID)
	notes.Put("/:id", middleware.Validate(new(schema.Note)), handler.UpdateNote)
	notes.Delete("/:id", handler.DeleteNote)
}
