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
		SELECT n.id, n.title,
			CASE
				WHEN length(n.content) > 100 THEN
					LEFT(n.content, 100 - POSITION(' ' IN REVERSE(LEFT(n.content, 100)))) || '...'
				ELSE n.content
			END as content,
			n.created_at,
			n.updated_at
		FROM notes n
		LEFT JOIN notes_users nu ON n.id = nu.note_id AND nu.user_id = $1
		WHERE n.user_id = $1 OR nu.user_id = $1
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
	query := `
		SELECT n.id, n.title, n.content, n.created_at, n.updated_at
		FROM notes n
		LEFT JOIN notes_users nu ON n.id = nu.note_id AND nu.user_id = $2
		WHERE n.id = $1 AND (user_id = $2 OR nu.user_id = $2)
	`
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
	query := `
		UPDATE notes n
		SET title = $1, content = $2
		FROM notes_users nu
		WHERE n.id = $3
			AND (n.user_id = $4 OR (nu.note_id = $3 AND nu.user_id = $4 AND nu.role = 'editor'))
			AND n.id = nu.note_id
	`
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
