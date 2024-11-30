package middleware

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/amiftachulh/notez-api/model"
	"github.com/gofiber/fiber/v2"
)

func ValidateBody[T interface {
	New() interface{}
	Validate() error
}](schema T) fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := schema.New()
		if err := c.BodyParser(body); err != nil {
			var syntaxErr *json.SyntaxError
			if errors.As(err, &syntaxErr) {
				return c.Status(fiber.StatusBadRequest).JSON(model.Response{
					Message: "Malformed JSON.",
				})
			}

			var unmarshalTypeErr *json.UnmarshalTypeError
			if errors.As(err, &unmarshalTypeErr) {
				errMsg := fmt.Sprintf(
					"Invalid value for field '%s'. Expected type '%s'.",
					unmarshalTypeErr.Field,
					unmarshalTypeErr.Type,
				)

				return c.Status(fiber.StatusBadRequest).JSON(model.Response{
					Message: "Invalid JSON type.",
					Error:   errMsg,
				})
			}
		}

		if err := body.(T).Validate(); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(model.Response{
				Message: "Validation failed.",
				Error:   err,
			})
		}

		c.Locals("body", body)

		return c.Next()
	}
}

func ValidateParams[T interface {
	New() interface{}
}](schema T) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := schema.New()
		if err := c.ParamsParser(params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Message: "Invalid parameter.",
				Error:   err.Error(),
			})
		}

		c.Locals("params", params)

		return c.Next()
	}
}

func ValidateQuery[T interface {
	New() interface{}
	Validate() error
}](schema T) fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := schema.New()
		c.QueryParser(query)

		if err := query.(T).Validate(); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(model.Response{
				Message: "Query validation failed.",
				Error:   err,
			})
		}

		c.Locals("query", query)

		return c.Next()
	}
}
