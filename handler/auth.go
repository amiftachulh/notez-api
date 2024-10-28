package handler

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"errors"
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
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func Login(c *fiber.Ctx) error {
	body := c.Locals("body").(schema.Login)

	var u user
	query := "SELECT * FROM users WHERE lower(email) = lower($1)"
	row := db.DB.QueryRow(query, body.Email)
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return c.Status(fiber.StatusUnauthorized).JSON(model.ErrResp{
				Message: "Invalid email or password.",
			})
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	match, err := argon2id.ComparePasswordAndHash(body.Password, u.Password)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if !match {
		return c.Status(fiber.StatusUnauthorized).JSON(model.ErrResp{
			Message: "Invalid email or password.",
		})
	}

	bytes := make([]byte, 15)
	rand.Read(bytes)
	sessionID := base32.StdEncoding.EncodeToString(bytes)

	query = "INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, now() + interval '7 day')"
	_, err = db.DB.Exec(query, sessionID, u.ID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(u)
}
