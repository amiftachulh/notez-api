package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/amiftachulh/notez-api/db"
	"github.com/amiftachulh/notez-api/model"
	"github.com/google/uuid"
)

func CreateSession(sessionID string, userID uuid.UUID, expiresAt time.Time) error {
	query := "INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)"
	if _, err := db.DB.Exec(query, sessionID, userID, expiresAt); err != nil {
		return err
	}
	return nil
}

func GetUserBySession(sessionID string) (*model.AuthUser, error) {
	var u model.AuthUser
	query := `
		SELECT u.id, u.name, u.email, u.role, u.created_at, u.updated_at, s.expires_at
		FROM sessions s
		JOIN users u
		ON s.user_id = u.id
		WHERE s.id = $1 AND s.expires_at > now()
	`
	if err := db.DB.
		QueryRow(query, sessionID).
		Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt, &u.ExpiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func DeleteSession(sessionID string) (bool, error) {
	query := "DELETE FROM sessions WHERE id = $1"
	result, err := db.DB.Exec(query, sessionID)
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
