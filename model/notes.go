package model

import (
	"github.com/invopop/validation"
)

type Note struct {
	Title   string  `json:"title"`
	Content *string `json:"content"`
}

func (c Note) Validate() error {
	return validation.ValidateStruct(
		&c,
		validation.Field(
			&c.Title,
			validation.Required.Error("Title is required."),
			validation.RuneLength(1, 300).Error("Title must be less than 300 characters."),
		),
		validation.Field(
			&c.Content,
			validation.When(
				c.Content != nil,
				validation.Required.Error("Content must be at least 1 character."),
				validation.RuneLength(1, 5000).Error("Content must be less than 5000 characters."),
			),
		),
	)
}
