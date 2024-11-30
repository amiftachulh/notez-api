package handler

import (
	"log"

	"github.com/amiftachulh/notez-api/model"
	"github.com/amiftachulh/notez-api/service"

	"github.com/gofiber/fiber/v2"
)

func CreateNote(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.AuthUser)
	body := c.Locals("body").(*model.NoteInput)

	if err := service.CreateNote(body.Title, body.Content, auth.ID); err != nil {
		log.Println("Error creating note:", err)
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Message: "Note created.",
	})
}

func GetNotes(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.AuthUser)
	query := c.Locals("query").(*model.NoteQuery)

	notes, total, err := service.GetNotes(auth.ID, query)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(model.PaginationResponse{
		Total: total,
		Items: notes,
	})
}

func GetNoteByID(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.AuthUser)
	id := c.Locals("params").(*model.NoteParams).ID

	note, err := service.GetNoteByID(id, auth.ID)
	if err != nil {
		log.Println("Error getting note by ID:", err)
		return fiber.ErrInternalServerError
	}
	if note == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Message: "Note not found.",
		})
	}

	return c.JSON(note)
}

func UpdateNoteByID(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.AuthUser)
	id := c.Locals("params").(*model.NoteParams).ID
	body := c.Locals("body").(*model.NoteInput)

	result, err := service.UpdateNoteByID(body, id, auth.ID)
	if err != nil {
		log.Println("Error updating note:", err)
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

func DeleteNoteByID(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.AuthUser)
	id := c.Locals("params").(*model.NoteParams).ID

	result, err := service.DeleteNoteByID(id, auth.ID)
	if err != nil {
		log.Println("Error deleting note:", err)
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
