package util

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/amiftachulh/notez-api/model"
)

func HandleJSONError(err error) model.Response {
	var jsonErr *json.UnmarshalTypeError
	if errors.As(err, &jsonErr) {
		return model.Response{
			Message: "Invalid input for field: " + jsonErr.Field,
			Error: fmt.Sprintf(
				"Expected type %s, but received %s.",
				jsonErr.Type.String(),
				jsonErr.Value,
			),
		}
	}
	return model.Response{
		Message: "Invalid JSON.",
	}
}
