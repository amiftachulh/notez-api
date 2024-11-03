package handler

import (
	"log"

	"github.com/amiftachulh/notez-api/model"
	"github.com/amiftachulh/notez-api/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateNoteInvitation(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.AuthUser)
	body := c.Locals("body").(*model.CreateNoteInvitation)

	if auth.Email == body.Email {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Message: "You can't invite yourself.",
		})
	}

	noteExists, err := repository.CheckNoteExists(body.NoteID, auth.ID)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	if !noteExists {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Message: "Note not found.",
		})
	}

	targetUserID, err := repository.GetUserIDByEmail(body.Email)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	if targetUserID == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Message: "User not found.",
		})
	}

	inviteExists, err := repository.CheckInviteExists(body.NoteID, *targetUserID)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	if inviteExists {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Message: "User already invited to note.",
		})
	}

	err = repository.CreateNoteInvitation(body.NoteID, auth.ID, *targetUserID, body.Role)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Message: "User invited to note.",
	})
}

func GetNoteInvitations(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.AuthUser)
	invitations, err := repository.GetNoteInvitations(auth.ID)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	return c.JSON(invitations)
}

func RespondNoteInvitation(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.AuthUser)
	body := c.Locals("body").(*model.RespondNoteInvitation)
	invitationID := c.Params("id")

	id, err := uuid.Parse(invitationID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Message: "Invalid invitation ID.",
		})
	}

	if !body.Accept {
		err := repository.DeclineInvitation(id, auth.ID)
		if err != nil {
			log.Println(err)
			return fiber.ErrInternalServerError
		}
		return c.JSON(model.Response{
			Message: "Invitation declined.",
		})
	}

	noteInvitation, err := repository.GetNoteInvitationByID(id, auth.ID)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	if noteInvitation == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Message: "Invitation not found.",
		})
	}

	err = repository.AcceptInvitation(noteInvitation.NoteID, auth.ID, noteInvitation.Role)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(model.Response{
		Message: "Invitation accepted.",
	})
}
