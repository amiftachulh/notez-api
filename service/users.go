package service

import (
	"database/sql"
	"errors"

	"github.com/amiftachulh/notez-api/db"
	"github.com/amiftachulh/notez-api/model"

	"github.com/google/uuid"
)

func CheckEmailExists(email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := db.DB.QueryRow(query, email).Scan(&exists)
	return exists, err
}

func CreateUser(email, password string) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	query := "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)"
	_, err = db.DB.Exec(query, id, email, password)
	return err
}

func GetUserByID(userID uuid.UUID) (*model.User, error) {
	var u model.User
	query := "SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE id = $1"
	err := db.DB.
		QueryRow(query, userID).
		Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func GetUserByEmail(email string) (*model.AuthUser, error) {
	var u model.AuthUser
	query := "SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE email = $1"
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

func UpdateUserInfo(userID uuid.UUID, body *model.UpdateUserInfo) (bool, error) {
	query := "UPDATE users SET name = $1 WHERE id = $2"
	result, err := db.DB.Exec(query, body.Name, userID)
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

func UpdateUserEmail(userID uuid.UUID, email string) (bool, error) {
	query := "UPDATE users SET email = $1 WHERE id = $2"
	result, err := db.DB.Exec(query, email, userID)
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

func UpdateUserPassword(userID uuid.UUID, hashedPassword string) (bool, error) {
	query := "UPDATE users SET password = $1 WHERE id = $2"
	result, err := db.DB.Exec(query, hashedPassword, userID)
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
