package handler

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"notez-api/db"
	"notez-api/model"
	"notez-api/schema"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Register(c *fiber.Ctx) error {
	body := c.Locals("body").(*schema.Register)

	var exist bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE lower(email) = lower($1))"
	err := db.DB.QueryRow(query, body.Email).Scan(&exist)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if exist {
		return c.Status(fiber.StatusConflict).JSON(model.ErrResp{
			Message: "Email already used.",
		})
	}

	hash, err := argon2id.CreateHash(body.Password, &argon2id.Params{
		Memory:      19456,
		Iterations:  2,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	})
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	query = "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)"
	_, err = db.DB.Exec(query, id, body.Email, hash)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(model.ErrResp{
		Message: "Register success.",
	})
}

type user struct {
	ID        uuid.UUID `json:"id"`
	Name      *string   `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func Login(c *fiber.Ctx) error {
	body := c.Locals("body").(*schema.Login)

	var u user
	query := "SELECT * FROM users WHERE lower(email) = lower($1)"
	err := db.DB.
		QueryRow(query, body.Email).
		Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusUnauthorized).JSON(model.ErrResp{
				Message: "Invalid email or password.",
			})
		}
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	match, err := argon2id.ComparePasswordAndHash(body.Password, u.Password)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if !match {
		return c.Status(fiber.StatusUnauthorized).JSON(model.ErrResp{
			Message: "Invalid email or password.",
		})
	}

	bytes := make([]byte, 15)
	rand.Read(bytes)
	sessionID := base64.RawURLEncoding.EncodeToString(bytes)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	u.ExpiresAt = expiresAt

	query = "INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)"
	_, err = db.DB.Exec(query, sessionID, u.ID, expiresAt)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "session"
	cookie.Value = sessionID
	cookie.Expires = expiresAt
	cookie.HTTPOnly = true
	cookie.Secure = true
	c.Cookie(cookie)

	return c.JSON(u)
}

func Logout(c *fiber.Ctx) error {
	sessionID := c.Cookies("session")
	if sessionID == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	query := "DELETE FROM sessions WHERE id = $1"
	result, err := db.DB.Exec(query, sessionID)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if rows == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "session"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.HTTPOnly = true
	cookie.Secure = true
	c.Cookie(cookie)

	return c.SendStatus(fiber.StatusNoContent)
}

func CheckAuth(c *fiber.Ctx) error {
	sessionID := c.Cookies("session")
	if sessionID == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var u user
	query := `
		SELECT u.id, u.name, u.email, u.role, u.created_at, u.updated_at, s.expires_at
		FROM sessions s
		JOIN users u
		ON s.user_id = u.id
		WHERE s.id = $1
	`
	err := db.DB.
		QueryRow(query, sessionID).
		Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt, &u.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(u)
}
