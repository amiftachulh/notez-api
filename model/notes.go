package model

import (
	"github.com/google/uuid"
	"github.com/invopop/validation"
)

type Note struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title,omitempty"`
	Content   *string   `json:"content"`
	CreatedAt string    `json:"created_at,omitempty"`
	UpdatedAt string    `json:"updated_at,omitempty"`
}

type NoteInput struct {
	Title   string  `json:"title"`
	Content *string `json:"content"`
}

func (c NoteInput) Validate() error {
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
