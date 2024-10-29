package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"notez-api/model"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

func Validate(v interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.BodyParser(v); err != nil {
			var jsonErr *json.UnmarshalTypeError
			if errors.As(err, &jsonErr) {
				return c.Status(fiber.StatusBadRequest).JSON(model.Response{
					Message: "Invalid input for field: " + jsonErr.Field,
					Error: fmt.Sprintf(
						"Expected type %s, but received %s.",
						jsonErr.Type.String(),
						jsonErr.Value,
					),
				})
			}
			return c.Status(fiber.StatusBadRequest).JSON(model.Response{
				Message: "Invalid JSON.",
				Error:   err,
			})
		}

		if validator, ok := v.(interface{ Validate() error }); ok {
			if err := validator.Validate(); err != nil {
				return c.Status(fiber.StatusUnprocessableEntity).JSON(model.Response{
					Message: "Validation error.",
					Error:   err,
				})
			}
		} else {
			t := reflect.TypeOf(validator)
			log.Println(fmt.Sprintf("%s doesn't have Validator() method.)", t.Name()))
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		c.Locals("body", v)
		return c.Next()
	}
}
