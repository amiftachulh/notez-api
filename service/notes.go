package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

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

func GetNotes(userID uuid.UUID, opts *model.NoteQuery) ([]model.NoteResponse, int, error) {
	notes := []model.NoteResponse{}

	queryBuilder := strings.Builder{}
	countQueryBuilder := strings.Builder{}

	queryBuilder.WriteString(
		"SELECT n.id, n.user_id, n.title, nu.role, n.created_at, n.updated_at FROM notes n LEFT JOIN notes_users nu ON n.id = nu.note_id AND nu.user_id = $1 WHERE (n.user_id = $1 OR nu.user_id = $1)",
	)
	countQueryBuilder.WriteString(
		"SELECT COUNT(*) FROM notes n LEFT JOIN notes_users nu ON n.id = nu.note_id AND nu.user_id = $1 WHERE (n.user_id = $1 OR nu.user_id = $1)",
	)
	params := []interface{}{userID}

	if opts.Query != "" {
		str := fmt.Sprintf(" AND n.title ILIKE $%d", len(params)+1)
		queryBuilder.WriteString(str)
		countQueryBuilder.WriteString(str)
		params = append(params, "%"+opts.Query+"%")
	}

	if opts.Role != "" {
		if opts.Role == "owner" {
			str := fmt.Sprintf(" AND n.user_id = $%d", len(params)+1)
			queryBuilder.WriteString(str)
			countQueryBuilder.WriteString(str)
			params = append(params, userID)
		} else {
			str := fmt.Sprintf(" AND nu.role = $%d", len(params)+1)
			queryBuilder.WriteString(str)
			countQueryBuilder.WriteString(str)
			params = append(params, opts.Role)
		}
	}

	queryBuilder.WriteString(
		fmt.Sprintf(
			" ORDER BY n.%s %s LIMIT %d OFFSET %d",
			opts.Sort,
			opts.Order,
			opts.PageSize,
			(opts.Page-1)*opts.PageSize,
		),
	)

	query := queryBuilder.String()
	rows, err := db.DB.Query(query, params...)
	if err != nil {
		log.Println("Error querying notes:", err)
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var n model.NoteResponse
		err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Role, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			log.Println("Error scanning note:", err)
		}
		notes = append(notes, n)
	}
	if err = rows.Err(); err != nil {
		log.Println("Error iterating notes:", err)
		return nil, 0, err
	}

	var total int
	countQuery := countQueryBuilder.String()
	if err = db.DB.QueryRow(countQuery, params...).Scan(&total); err != nil {
		log.Println("Error counting notes:", err)
		return nil, 0, err
	}

	return notes, total, nil
}

func GetNoteByID(noteID, userID uuid.UUID) (*model.NoteDetail, error) {
	var n model.NoteDetail
	query := `
		SELECT n.id, n.title, n.content, nu.role, u.id AS owner_id, u.email, u.name, n.created_at, n.updated_at
		FROM notes n
		JOIN users u ON n.user_id = u.id
		LEFT JOIN notes_users nu ON n.id = nu.note_id AND nu.user_id = $2
		WHERE n.id = $1 AND (n.user_id = $2 OR nu.user_id = $2)
	`
	err := db.DB.
		QueryRow(query, noteID, userID).
		Scan(&n.ID, &n.Title, &n.Content, &n.Role, &n.Owner.ID, &n.Owner.Email, &n.Owner.Name, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	members := []model.NoteMember{}
	query = `
		SELECT u.id, u.email, u.name, nu.role, nu.created_at
		FROM notes_users nu
		JOIN users u ON nu.user_id = u.id
		WHERE nu.note_id = $1
	`
	rows, err := db.DB.Query(query, noteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var m model.NoteMember
		if err := rows.Scan(&m.ID, &m.Email, &m.Name, &m.Role, &m.CreatedAt); err != nil {
			log.Println(err)
		}
		members = append(members, m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	n.Members = members
	return &n, nil
}

func CheckNoteExists(id, userID uuid.UUID) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM notes WHERE id = $1 AND user_id = $2)"
	err := db.DB.QueryRow(query, id, userID).Scan(&exists)
	return exists, err
}

func UpdateNoteByID(body *model.NoteInput, noteID uuid.UUID, userID uuid.UUID) (bool, error) {
	query := `
		WITH notes_to_update AS (
			SELECT n.id
			FROM notes n
			LEFT JOIN notes_users nu ON nu.note_id = n.id
			WHERE n.id = $3
				AND (
					n.user_id = $4
					OR (
					    nu.user_id = $4
						AND nu.role = 'editor'
					)
				)
		)
		UPDATE notes
		SET title = $1, content = $2
		WHERE id IN (SELECT id FROM notes_to_update)
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

func DeleteNoteByID(id, userID uuid.UUID) (bool, error) {
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
