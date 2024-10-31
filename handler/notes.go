package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"notez-api/db"
	"notez-api/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateNote(c *fiber.Ctx) error {
	body := c.Locals("body").(*model.Note)
	user := c.Locals("user").(model.User)

	id, err := uuid.NewV7()
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	query := "INSERT INTO notes (id, title, content, user_id) VALUES ($1, $2, $3, $4)"
	_, err = db.DB.Exec(query, id, body.Title, body.Content, user.ID)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Message: "Note created.",
	})
}

func GetNotes(c *fiber.Ctx) error {
	user := c.Locals("user").(model.User)

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
	rows, err := db.DB.Query(query, user.ID)
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
	user := c.Locals("user").(model.User)
	noteID := c.Params("id")

	var n note
	query := "SELECT id, title, content, created_at, updated_at FROM notes WHERE id = $1 AND user_id = $2"
	err := db.DB.
		QueryRow(query, noteID, user.ID).
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
	user := c.Locals("user").(model.User)
	noteID := c.Params("id")

	query := "UPDATE notes SET title = $1, content = $2 WHERE id = $3 AND user_id = $4"
	result, err := db.DB.Exec(query, body.Title, body.Content, noteID, user.ID)
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
	user := c.Locals("user").(model.User)
	noteID := c.Params("id")

	query := "DELETE FROM notes WHERE id = $1 AND user_id = $2"
	result, err := db.DB.Exec(query, noteID, user.ID)
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

func InviteUserToNote(c *fiber.Ctx) error {
	user := c.Locals("user").(model.User)
	body := c.Locals("body").(*model.NoteInvite)
	noteID := c.Params("id")

	var noteExist bool
	query := "SELECT EXISTS(SELECT 1 FROM notes WHERE id = $1 AND user_id = $2)"
	err := db.DB.QueryRow(query, noteID, user.ID).Scan(&noteExist)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	if !noteExist {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Message: "Note not found.",
		})
	}

	var targetUserID uuid.UUID
	query = "SELECT id FROM users WHERE email = $1"
	if err = db.DB.QueryRow(query, body.Email).Scan(&targetUserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(model.Response{
				Message: fmt.Sprintf("User with email %s not found.", body.Email),
			})
		}
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	var inviteExist bool
	query = "SELECT EXISTS(SELECT 1 FROM note_invites WHERE note_id = $1 AND user_id = $2)"
	err = db.DB.QueryRow(query, noteID, targetUserID).Scan(&inviteExist)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	if inviteExist {
		return c.Status(fiber.StatusConflict).JSON(model.Response{
			Message: "User already invited.",
		})
	}

	query = "INSERT INTO note_invites (note_id, user_id) VALUES ($1, $2)"
	_, err = db.DB.Exec(query, noteID, targetUserID)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(model.Response{
		Message: "User invited to note.",
	})
}
