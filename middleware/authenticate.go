package middleware

import (
	"database/sql"
	"errors"
	"log"
	"notez-api/db"
	"notez-api/schema"

	"github.com/gofiber/fiber/v2"
)

func Authenticate(c *fiber.Ctx) error {
	sessionID := c.Cookies("session")
	if sessionID == "" {
		return fiber.ErrUnauthorized
	}

	var u schema.User
	query := `
		SELECT u.id, u.name, u.email, u.role, u.created_at, u.updated_at, s.expires_at
		FROM sessions s
		JOIN users u
		ON s.user_id = u.id
		WHERE s.id = $1 AND s.expires_at > now()
	`
	err := db.DB.
		QueryRow(query, sessionID).
		Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt, &u.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fiber.ErrUnauthorized
		}
		log.Println(err)
		return err
	}

	c.Locals("user", u)
	return c.Next()
}