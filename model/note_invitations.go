package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

type CreateNoteInvitation struct {
	Email  string    `json:"email"`
	NoteID uuid.UUID `json:"note_id"`
	Role   string    `json:"role"`
}

func (i CreateNoteInvitation) New() interface{} {
	return &CreateNoteInvitation{}
}

func (i CreateNoteInvitation) Validate() error {
	return validation.ValidateStruct(
		&i,
		validation.Field(
			&i.Email,
			validation.Required.Error("Email is required."),
			is.Email.Error("Email is invalid."),
		),
		validation.Field(
			&i.Role,
			validation.In("editor", "viewer").Error("Role must be either 'editor' or 'viewer'."),
		),
	)
}

type NoteInvitationParams struct {
	ID uuid.UUID `param:"id"`
}

func (p NoteInvitationParams) New() interface{} {
	return &NoteInvitationParams{}
}

type NoteInvitation struct {
	ID     uuid.UUID `json:"id"`
	NoteID uuid.UUID `json:"note_id"`
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
}

type noteWithoutUserID struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type NoteInvitationResponse struct {
	ID        uuid.UUID         `json:"id"`
	Note      noteWithoutUserID `json:"note"`
	Inviter   User              `json:"inviter"`
	Role      string            `json:"role"`
	CreatedAt time.Time         `json:"created_at"`
}

type RespondNoteInvitation struct {
	Accept bool `json:"accept"`
}
