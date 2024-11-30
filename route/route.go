package route

import (
	"github.com/amiftachulh/notez-api/handler"
	"github.com/amiftachulh/notez-api/middleware"
	"github.com/amiftachulh/notez-api/model"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	v1 := app.Group("/v1")

	auth := v1.Group("/auth")
	auth.Post("/register", middleware.ValidateBody(&model.Register{}), handler.Register)
	auth.Post("/login", middleware.ValidateBody(&model.Login{}), handler.Login)
	auth.Post("/logout", handler.Logout)
	auth.Get("/check", handler.CheckAuth)

	protected := v1.Group("/").Use(middleware.Authenticate)

	profile := protected.Group("/profile")
	profile.Patch("/", middleware.ValidateBody(&model.UpdateUserInfo{}), handler.UpdateUserInfo)
	profile.Patch(
		"/email",
		middleware.ValidateBody(&model.UpdateUserEmail{}),
		handler.UpdateUserEmail,
	)
	profile.Patch(
		"/password",
		middleware.ValidateBody(&model.UpdateUserPassword{}),
		handler.UpdateUserPassword,
	)

	notes := protected.Group("/notes")
	notes.Post("/", middleware.ValidateBody(&model.NoteInput{}), handler.CreateNote)
	notes.Get("/", middleware.ValidateQuery(&model.NoteQuery{}), handler.GetNotes)
	notes.Get("/:id", middleware.ValidateParams(&model.NoteParams{}), handler.GetNoteByID)
	notes.Put(
		"/:id",
		middleware.ValidateParams(&model.NoteParams{}),
		middleware.ValidateBody(&model.NoteInput{}),
		handler.UpdateNoteByID,
	)
	notes.Delete("/:id", middleware.ValidateParams(&model.NoteParams{}), handler.DeleteNoteByID)
	notes.Patch(
		"/:id/members/:memberID",
		middleware.ValidateParams(&model.NoteMemberParams{}),
		middleware.ValidateBody(&model.UpdateNoteMemberRole{}),
		handler.UpdateNoteMemberRole,
	)
	notes.Delete(
		"/:id/members/:memberID",
		middleware.ValidateParams(&model.NoteMemberParams{}),
		handler.RemoveNoteMember,
	)

	noteInvitation := protected.Group("/note-invitations")
	noteInvitation.Post(
		"/",
		middleware.ValidateBody(&model.CreateNoteInvitation{}),
		handler.CreateNoteInvitation,
	)
	noteInvitation.Get("/", handler.GetNoteInvitations)
	noteInvitation.Patch(
		"/:id",
		middleware.ValidateParams(&model.NoteInvitationParams{}),
		handler.RespondNoteInvitation,
	)
}
