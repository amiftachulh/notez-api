package handler

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"time"

	"github.com/amiftachulh/notez-api/model"
	"github.com/amiftachulh/notez-api/service"

	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	body := c.Locals("body").(*model.Register)

	exists, err := service.CheckEmailExists(body.Email)
	if err != nil {
		log.Println("Error checking email exists:", err)
		return fiber.ErrInternalServerError
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(model.Response{
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
		log.Println("Error creating hash:", err)
		return fiber.ErrInternalServerError
	}

	err = service.RegisterUser(body.Email, hash)
	if err != nil {
		log.Println("Error registering user:", err)
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Message: "Register success.",
	})
}

func Login(c *fiber.Ctx) error {
	body := c.Locals("body").(*model.Login)

	user, err := service.GetUserByEmail(body.Email)
	if err != nil {
		log.Println("Error getting user by email:", err)
		return fiber.ErrInternalServerError
	}
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
			Message: "Invalid email or password.",
		})
	}

	match, err := argon2id.ComparePasswordAndHash(body.Password, user.Password)
	if err != nil {
		log.Println("Error comparing password and hash:", err)
		return fiber.ErrInternalServerError
	}
	if !match {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
			Message: "Invalid email or password.",
		})
	}

	bytes := make([]byte, 15)
	rand.Read(bytes)
	sessionID := base64.RawURLEncoding.EncodeToString(bytes)
	expiresAt := time.Now().Add(7 * 24 * time.Hour).Truncate(time.Microsecond)
	user.ExpiresAt = expiresAt

	if err = service.CreateSession(sessionID, user.ID, expiresAt); err != nil {
		log.Println("Error creating session:", err)
		return fiber.ErrInternalServerError
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "session"
	cookie.Value = sessionID
	cookie.Expires = expiresAt
	cookie.HTTPOnly = true
	cookie.Secure = true
	c.Cookie(cookie)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	sessionID := c.Cookies("session")
	if sessionID == "" {
		return fiber.ErrUnauthorized
	}

	result, err := service.DeleteSession(sessionID)
	if err != nil {
		log.Println("Error deleting session:", err)
		return fiber.ErrInternalServerError
	}
	if !result {
		return fiber.ErrUnauthorized
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "session"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.HTTPOnly = true
	cookie.Secure = true
	c.Cookie(cookie)

	return c.JSON(model.Response{
		Message: "Logout success.",
	})
}

func CheckAuth(c *fiber.Ctx) error {
	sessionID := c.Cookies("session")
	if sessionID == "" {
		return fiber.ErrUnauthorized
	}

	user, err := service.GetUserBySession(sessionID)
	if err != nil {
		log.Println("Error checking auth:", err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(user)
}
