package route

import (
	"notez-api/handler"
	"notez-api/middleware"
	"notez-api/model"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	v1 := app.Group("/v1")

	auth := v1.Group("/auth")
	auth.Post("/register", middleware.Validate(new(model.Register)), handler.Register)
	auth.Post("/login", middleware.Validate(new(model.Login)), handler.Login)
	auth.Post("/logout", handler.Logout)
	auth.Get("/check", handler.CheckAuth)

	notes := v1.Group("/notes").Use(middleware.Authenticate)
	notes.Post("/", middleware.Validate(new(model.Note)), handler.CreateNote)
	notes.Get("/", handler.GetNotes)
	notes.Get("/:id", handler.GetNoteByID)
	notes.Put("/:id", middleware.Validate(new(model.Note)), handler.UpdateNote)
	notes.Delete("/:id", handler.DeleteNote)
	notes.Post("/:id/invite", middleware.Validate(new(model.NoteInvite)), handler.InviteUserToNote)
}
