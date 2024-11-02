package handler

import (
	"database/sql"
	"errors"
	"log"
	"notez-api/db"
	"notez-api/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateNote(c *fiber.Ctx) error {
	body := c.Locals("body").(*model.Note)
	auth := c.Locals("auth").(model.User)

	id, err := uuid.NewV7()
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	query := "INSERT INTO notes (id, title, content, user_id) VALUES ($1, $2, $3, $4)"
	_, err = db.DB.Exec(query, id, body.Title, body.Content, auth.ID)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Message: "Note created.",
	})
}

func GetNotes(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.User)

	notes := []note{}
	query := `
		SELECT id, title,
			CASE
				WHEN length(content) > 100 THEN
					LEFT(content, 100 - POSITION(' ' IN REVERSE(LEFT(content, 100)))) || '...'
				ELSE content
			END as content,
			created_at,
			updated_at
		FROM notes
		WHERE user_id = $1
		LIMIT 10
	`
	rows, err := db.DB.Query(query, auth.ID)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	defer rows.Close()
	for rows.Next() {
		var n note
		err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			log.Println(err)
		}
		notes = append(notes, n)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(notes)
}

func GetNoteByID(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.User)
	noteID := c.Params("id")

	if _, err := uuid.Parse(noteID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Message: "Invalid note ID.",
		})
	}

	var n note
	query := "SELECT id, title, content, created_at, updated_at FROM notes WHERE id = $1 AND user_id = $2"
	err := db.DB.
		QueryRow(query, noteID, auth.ID).
		Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(model.Response{
				Message: "Note not found.",
			})
		}
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(n)
}

func UpdateNote(c *fiber.Ctx) error {
	body := c.Locals("body").(*model.Note)
	auth := c.Locals("auth").(model.User)
	noteID := c.Params("id")

	if _, err := uuid.Parse(noteID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Message: "Invalid note ID.",
		})
	}

	query := "UPDATE notes SET title = $1, content = $2 WHERE id = $3 AND user_id = $4"
	result, err := db.DB.Exec(query, body.Title, body.Content, noteID, auth.ID)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	if rows == 0 {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Message: "Note not found.",
		})
	}

	return c.JSON(model.Response{
		Message: "Note updated.",
	})
}

func DeleteNote(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.User)
	noteID := c.Params("id")

	if _, err := uuid.Parse(noteID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Message: "Invalid note ID.",
		})
	}

	query := "DELETE FROM notes WHERE id = $1 AND user_id = $2"
	result, err := db.DB.Exec(query, noteID, auth.ID)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	if rows == 0 {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Message: "Note not found.",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
