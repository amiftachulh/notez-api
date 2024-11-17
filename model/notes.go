package model

import (
	"github.com/google/uuid"
	"github.com/invopop/validation"
)

const NOTE_MAX_CONTENT_BYTES = 5 * 1024 * 1024 // 5 MB

type Note struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Content   *string   `json:"content,omitempty"`
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
				validation.Length(1, NOTE_MAX_CONTENT_BYTES).
					Error("Content size must be less than 5 MB."),
			),
		),
	)
}

type NoteResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Content   *string   `json:"content"`
	Role      *string   `json:"role,omitempty"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}
