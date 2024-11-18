package util

import (
	"fmt"
	"strings"

	"github.com/amiftachulh/notez-api/model"
)

func HandleQueryError(err error) model.Response {
	if strings.Contains(err.Error(), "schema: error converting value") {
		field := extractFieldFromError(err.Error())
		return model.Response{
			Message: fmt.Sprintf(
				"Invalid value for query parameter '%s'.",
				field,
			),
		}
	}

	return model.Response{
		Message: "Invalid query parameters.",
		Error:   err.Error(),
	}
}

func extractFieldFromError(errMessage string) string {
	start := strings.Index(errMessage, "\"")
	end := strings.LastIndex(errMessage, "\"")
	if start != -1 && end != -1 && start < end {
		return errMessage[start+1 : end]
	}
	return "unknown"
}
