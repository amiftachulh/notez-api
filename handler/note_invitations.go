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

func CreateNoteInvitation(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.User)
	body := c.Locals("body").(*model.NoteInvitation)

	if auth.Email == body.Email {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Message: "You can't invite yourself.",
		})
	}

	var noteExist bool
	query := "SELECT EXISTS(SELECT 1 FROM notes WHERE id = $1 AND user_id = $2)"
	err := db.DB.QueryRow(query, body.NoteID, auth.ID).Scan(&noteExist)
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
				Message: fmt.Sprintf("User with email '%s' not found.", body.Email),
			})
		}
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	var inviteExist bool
	query = "SELECT EXISTS(SELECT 1 FROM note_invitations WHERE note_id = $1 AND user_id = $2)"
	err = db.DB.QueryRow(query, body.NoteID, targetUserID).Scan(&inviteExist)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	if inviteExist {
		return c.Status(fiber.StatusConflict).JSON(model.Response{
			Message: "User already invited.",
		})
	}

	id, err := uuid.NewV7()
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	query = "INSERT INTO note_invitations (id, note_id, user_id, inviter_id, role) VALUES ($1, $2, $3, $4, $5)"
	if _, err = db.DB.Exec(query, id, body.NoteID, targetUserID, auth.ID, body.Role); err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Message: "User invited to note.",
	})
}

func GetNoteInvitations(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.User)

	query := `
		SELECT ni.id, n.id, n.title, i.id, i.email, i.name, ni.role, ni.created_at
		FROM note_invitations ni
		JOIN notes n ON ni.note_id = n.id
		JOIN users i ON ni.inviter_id = i.id
		WHERE ni.user_id = $1
	`
	rows, err := db.DB.Query(query, auth.ID)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	invitations := []noteInvitationResponse{}

	defer rows.Close()
	for rows.Next() {
		var ni noteInvitationResponse
		if err := rows.Scan(
			&ni.ID,
			&ni.Note.ID,
			&ni.Note.Title,
			&ni.Inviter.ID,
			&ni.Inviter.Email,
			&ni.Inviter.Name,
			&ni.Role,
			&ni.CreatedAt,
		); err != nil {
			log.Println(err)
		}
		invitations = append(invitations, ni)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(invitations)
}

func RespondNoteInvitation(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.User)
	body := c.Locals("body").(*model.RespondNoteInvitation)
	invitationID := c.Params("id")

	if _, err := uuid.Parse(invitationID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Message: "Invalid invitation ID.",
		})
	}

	var query string
	if !body.Accept {
		query = "DELETE FROM note_invitations WHERE id = $1 AND user_id = $2"
		if _, err := db.DB.Exec(query, invitationID, auth.ID); err != nil {
			log.Println(err)
			return fiber.ErrInternalServerError
		}

		return c.JSON(model.Response{
			Message: "Invitation declined.",
		})
	}

	var noteID uuid.UUID
	var role string
	query = "SELECT note_id, role FROM note_invitations WHERE id = $1 AND user_id = $2"
	if err := db.DB.QueryRow(query, invitationID, auth.ID).Scan(&noteID, &role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(model.Response{
				Message: "Invitation not found.",
			})
		}
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	if err := acceptInvitation(noteID, auth.ID, role); err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(model.Response{
		Message: "Invitation accepted.",
	})
}

func acceptInvitation(noteID uuid.UUID, userID uuid.UUID, role string) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := "DELETE FROM note_invitations WHERE note_id = $1 AND user_id = $2"
	if _, err = tx.Exec(query, noteID, userID); err != nil {
		return err
	}

	query = "INSERT INTO notes_users (note_id, user_id, role) VALUES ($1, $2, $3)"
	if _, err = tx.Exec(query, noteID, userID, role); err != nil {
		return err
	}

	return tx.Commit()
}
