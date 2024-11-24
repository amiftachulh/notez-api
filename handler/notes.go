package handler

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/amiftachulh/notez-api/model"
	"github.com/amiftachulh/notez-api/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateNote(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.AuthUser)

	body := new(model.NoteInput)
	if err := c.BodyParser(body); err != nil {
		var syntaxError *json.SyntaxError
		if errors.As(err, &syntaxError) {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Message: invalidJSON,
			})
		}
	}

	if err := body.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(model.Response{
			Message: validationErr,
			Error:   err,
		})
	}

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

	query := &model.NoteQuery{
		Page:     1,
		PageSize: 10,
		Sort:     "id",
		Order:    "asc",
	}
	c.QueryParser(query)

	if err := query.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(model.Response{
			Message: queryValidationErr,
			Error:   err,
		})
	}

	notes, total, err := service.GetNotes(auth.ID, query)
	if err != nil {
		log.Println("Error getting notes:", err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(model.PaginationResponse{
		Total: total,
		Items: notes,
	})
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

	noteID := c.Params("id")
	id, err := uuid.Parse(noteID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Message: "Invalid note ID.",
		})
	}

	body := new(model.NoteInput)
	if err := c.BodyParser(body); err != nil {
		var syntaxError *json.SyntaxError
		if errors.As(err, &syntaxError) {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Message: invalidJSON,
			})
		}
	}

	if err := body.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(model.Response{
			Message: validationErr,
			Error:   err,
		})
	}

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

	noteID := c.Params("id")
	id, err := uuid.Parse(noteID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Message: "Invalid note ID.",
		})
	}

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
