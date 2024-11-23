package route

import (
	"github.com/amiftachulh/notez-api/handler"
	"github.com/amiftachulh/notez-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	v1 := app.Group("/v1")

	auth := v1.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)
	auth.Post("/logout", handler.Logout)
	auth.Get("/check", handler.CheckAuth)

	protected := v1.Group("/").Use(middleware.Authenticate)

	notes := protected.Group("/notes")
	notes.Post("/", handler.CreateNote)
	notes.Get("/", handler.GetNotes)
	notes.Get("/:id", handler.GetNoteByID)
	notes.Put("/:id", handler.UpdateNoteByID)
	notes.Delete("/:id", handler.DeleteNoteByID)
	notes.Patch("/:id/members/:memberID", handler.UpdateNoteMemberRole)
	notes.Delete("/:id/members/:memberID", handler.RemoveNoteMember)

	noteInvitation := protected.Group("/note-invitations")
	noteInvitation.Post("/", handler.CreateNoteInvitation)
	noteInvitation.Get("/", handler.GetNoteInvitations)
	noteInvitation.Patch("/:id", handler.RespondNoteInvitation)
}
