package handler

import (
	"errors"

	"github.com/amiftachulh/notez-api/model"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	return c.Status(code).JSON(model.Response{
		Message: err.Error(),
	})
}
