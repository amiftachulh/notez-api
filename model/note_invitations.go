package model

import (
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

type NoteInvitation struct {
	Email  string `json:"email"`
	NoteID string `json:"note_id"`
	Role   string `json:"role"`
}

func (i NoteInvitation) Validate() error {
	return validation.ValidateStruct(
		&i,
		validation.Field(
			&i.Email,
			validation.Required.Error("Email is required."),
			is.Email.Error("Email is invalid."),
		),
		validation.Field(
			&i.NoteID,
			validation.Required.Error("Note ID is required."),
			is.UUID.Error("Invalid note ID."),
		),
		validation.Field(
			&i.Role,
			validation.In("editor", "viewer").Error("Role must be either editor or viewer."),
		),
	)
}

type RespondNoteInvitation struct {
	Accept bool `json:"accept"`
}
