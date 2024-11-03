package handler

import (
	"log"

	"github.com/amiftachulh/notez-api/model"
	"github.com/amiftachulh/notez-api/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateNote(c *fiber.Ctx) error {
	body := c.Locals("body").(*model.NoteInput)
	auth := c.Locals("auth").(model.AuthUser)

	if err := repository.CreateNote(body.Title, body.Content, auth.ID); err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Message: "Note created.",
	})
}

func GetNotes(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.AuthUser)

	notes, err := repository.GetNotes(auth.ID)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(notes)
}

func GetNoteByID(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.AuthUser)
	noteID := c.Params("id")

	id, err := uuid.Parse(noteID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Message: "Invalid note ID.",
		})
	}

	n, err := repository.GetNoteByID(id, auth.ID)
	return c.JSON(n)
}

func UpdateNote(c *fiber.Ctx) error {
	body := c.Locals("body").(*model.NoteInput)
	auth := c.Locals("auth").(model.AuthUser)
	noteID := c.Params("id")

	id, err := uuid.Parse(noteID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Message: "Invalid note ID.",
		})
	}

	result, err := repository.UpdateNoteByID(body, id, auth.ID)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	if !result {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Message: "Note not found.",
		})
	}

	return c.JSON(model.Response{
		Message: "Note updated.",
	})
}

func DeleteNote(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.AuthUser)
	noteID := c.Params("id")

	id, err := uuid.Parse(noteID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Message: "Invalid note ID.",
		})
	}

	result, err := repository.DeleteNoteByID(id, auth.ID)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	if !result {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Message: "Note not found.",
		})
	}

	return c.JSON(model.Response{
		Message: "Note deleted.",
	})
}
