package service

import (
	"database/sql"
	"errors"
	"log"

	"github.com/amiftachulh/notez-api/db"
	"github.com/amiftachulh/notez-api/model"
	"github.com/google/uuid"
)

func CheckInviteExists(noteID uuid.UUID, targetUserID uuid.UUID) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM note_invitations WHERE note_id = $1 AND target_user_id = $2)"
	err := db.DB.QueryRow(query, noteID, targetUserID).Scan(&exists)
	return exists, err
}

func CreateNoteInvitation(noteID uuid.UUID, targetUserID uuid.UUID, inviterID uuid.UUID, role string) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	query := "INSERT INTO note_invitations (id, note_id, user_id, inviter_id, role) VALUES ($1, $2, $3, $4, $5)"
	_, err = db.DB.Exec(query, id, noteID, targetUserID, inviterID, role)
	return err
}

func GetNoteInvitations(userID uuid.UUID) ([]model.NoteInvitationResponse, error) {
	query := `
		SELECT ni.id, n.id, n.title, i.id, i.email, i.name, ni.role, ni.created_at
		FROM note_invitations ni
		JOIN notes n ON ni.note_id = n.id
		JOIN users i ON ni.inviter_id = i.id
		WHERE ni.user_id = $1
	`
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}

	invitations := []model.NoteInvitationResponse{}

	defer rows.Close()
	for rows.Next() {
		var ni model.NoteInvitationResponse
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
		return nil, err
	}
	return invitations, nil
}

func GetNoteInvitationByID(invitationID uuid.UUID, userID uuid.UUID) (*model.NoteInvitation, error) {
	var ni model.NoteInvitation
	query := "SELECT note_id, role FROM note_invitations WHERE id = $1 AND user_id = $2"
	if err := db.DB.QueryRow(query, invitationID, userID).Scan(&ni.ID, &ni.Role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &ni, nil
}

func DeclineInvitation(invitationID uuid.UUID, userID uuid.UUID) error {
	query := "DELETE FROM note_invitations WHERE id = $1 AND user_id = $2"
	_, err := db.DB.Exec(query, invitationID, userID)
	return err
}

func AcceptInvitation(noteID uuid.UUID, userID uuid.UUID, role string) error {
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
