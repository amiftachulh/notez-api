package service

import (
	"github.com/amiftachulh/notez-api/db"
	"github.com/google/uuid"
)

func CheckIsNoteOwner(noteID, userID uuid.UUID) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM notes WHERE id = $1 AND user_id = $2)"
	err := db.DB.QueryRow(query, noteID, userID).Scan(&exists)
	return exists, err
}

func UpdateNoteMemberRole(noteID, memberID uuid.UUID, role string) (bool, error) {
	query := "UPDATE notes_users SET role = $1 WHERE note_id = $2 AND user_id = $3"
	result, err := db.DB.Exec(query, role, noteID, memberID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func RemoveNoteMember(noteID, memberID uuid.UUID) (bool, error) {
	query := "DELETE FROM notes_users WHERE note_id = $1 AND user_id = $2"
	result, err := db.DB.Exec(query, noteID, memberID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected == 0 {
		return false, nil
	}

	return true, nil
}
