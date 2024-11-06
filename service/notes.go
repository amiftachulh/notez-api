package service

import (
	"database/sql"
	"errors"
	"log"

	"github.com/amiftachulh/notez-api/db"
	"github.com/amiftachulh/notez-api/model"

	"github.com/google/uuid"
)

func CreateNote(title string, content *string, userID uuid.UUID) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	query := "INSERT INTO notes (id, title, content, user_id) VALUES ($1, $2, $3, $4)"
	_, err = db.DB.Exec(query, id, title, content, userID)
	return err
}

func GetNotes(userID uuid.UUID) ([]model.Note, error) {
	notes := []model.Note{}
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
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var n model.Note
		err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			log.Println(err)
		}
		notes = append(notes, n)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return notes, nil
}

func GetNoteByID(noteID uuid.UUID, userID uuid.UUID) (*model.Note, error) {
	var n model.Note
	query := "SELECT id, title, content, created_at, updated_at FROM notes WHERE id = $1 AND user_id = $2"
	err := db.DB.
		QueryRow(query, noteID, userID).
		Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &n, nil
}

func CheckNoteExists(id uuid.UUID, userID uuid.UUID) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM notes WHERE id = $1 AND user_id = $2)"
	err := db.DB.QueryRow(query, id, userID).Scan(&exists)
	return exists, err
}

func UpdateNoteByID(body *model.NoteInput, noteID uuid.UUID, userID uuid.UUID) (bool, error) {
	query := "UPDATE notes SET title = $1, content = $2 WHERE id = $3 AND user_id = $4"
	result, err := db.DB.Exec(query, body.Title, body.Content, noteID, userID)
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

func DeleteNoteByID(id uuid.UUID, userID uuid.UUID) (bool, error) {
	query := "DELETE FROM notes WHERE id = $1 AND user_id = $2"
	result, err := db.DB.Exec(query, id, userID)
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
