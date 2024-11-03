package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/amiftachulh/notez-api/db"
	"github.com/amiftachulh/notez-api/model"

	"github.com/google/uuid"
)

func CheckEmailExists(email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	if err := db.DB.QueryRow(query, email).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

func RegisterUser(email, password string) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	query := "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)"
	_, err = db.DB.Exec(query, id, email, password)
	return err
}

func GetUserByEmail(email string) (*model.AuthUser, error) {
	var u model.AuthUser
	query := "SELECT * FROM users WHERE email = $1"
	if err := db.DB.
		QueryRow(query, email).
		Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func GetUserIDByEmail(email string) (*uuid.UUID, error) {
	var id uuid.UUID
	query := "SELECT id FROM users WHERE email = $1"
	if err := db.DB.QueryRow(query, email).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &id, nil
}

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
