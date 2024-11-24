package handler

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/alexedwards/argon2id"
	"github.com/amiftachulh/notez-api/model"
	"github.com/amiftachulh/notez-api/service"
	"github.com/gofiber/fiber/v2"
)

func UpdateUserInfo(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.AuthUser)

	body := new(model.UpdateUserInfo)
	if err := c.BodyParser(body); err != nil {
		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Message: invalidJSON,
			})
		}
	}

	if err := body.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(model.Response{
			Message: validationErr,
			Error:   err,
		})
	}

	result, err := service.UpdateUserInfo(auth.ID, body)
	if err != nil {
		log.Println("Error updating user info:", err)
		return fiber.ErrInternalServerError
	}
	if !result {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Message: userNotFound,
		})
	}

	return c.JSON(model.Response{
		Message: "User info updated.",
	})
}

func UpdateUserEmail(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.AuthUser)

	body := new(model.UpdateUserEmail)
	if err := c.BodyParser(body); err != nil {
		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Message: invalidJSON,
			})
		}
	}

	if err := body.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(model.Response{
			Message: "Validation error.",
			Error:   err,
		})
	}

	if body.Email == auth.Email {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Message: "Email still the same.",
		})
	}

	exists, err := service.CheckEmailExists(body.Email)
	if err != nil {
		log.Println("Error checking email exists:", err)
		return fiber.ErrInternalServerError
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(model.Response{
			Message: emailUsed,
		})
	}

	result, err := service.UpdateUserEmail(auth.ID, body.Email)
	if err != nil {
		log.Println("Error updating user email:", err)
		return fiber.ErrInternalServerError
	}
	if !result {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Message: userNotFound,
		})
	}

	return c.JSON(model.Response{
		Message: "User email updated.",
	})
}

func UpdateUserPassword(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.AuthUser)

	body := new(model.UpdateUserPassword)
	if err := c.BodyParser(body); err != nil {
		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Message: invalidJSON,
			})
		}
	}

	if err := body.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(model.Response{
			Message: validationErr,
			Error:   err,
		})
	}

	user, err := service.GetUserByID(auth.ID)
	if err != nil {
		log.Println("Error getting user by ID:", err)
		return fiber.ErrInternalServerError
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Message: userNotFound,
		})
	}

	match, err := argon2id.ComparePasswordAndHash(body.CurrentPassword, user.Password)
	if err != nil {
		log.Println("Error comparing password and hash:", err)
		return fiber.ErrInternalServerError
	}
	if !match {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(model.Response{
			Message: updatePasswordFail,
			Error: map[string]string{
				"current_password": "Current password is incorrect.",
			},
		})
	}

	if body.CurrentPassword == body.Password {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(model.Response{
			Message: updatePasswordFail,
			Error: map[string]string{
				"password": "New password can't be the same as the current password.",
			},
		})
	}

	hash, err := hashPassword(body.Password)
	if err != nil {
		log.Println("Error creating hash:", err)
		return fiber.ErrInternalServerError
	}

	result, err := service.UpdateUserPassword(auth.ID, hash)
	if err != nil {
		log.Println("Error updating user password:", err)
		return fiber.ErrInternalServerError
	}
	if !result {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Message: userNotFound,
		})
	}

	return c.JSON(model.Response{
		Message: "User password updated.",
	})
}
