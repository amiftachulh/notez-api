package model

import (
	"github.com/google/uuid"
	"github.com/invopop/validation"
)

type UpdateNoteMemberRole struct {
	Role string `json:"role"`
}

func (r UpdateNoteMemberRole) New() interface{} {
	return &UpdateNoteMemberRole{}
}

func (r UpdateNoteMemberRole) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(
			&r.Role,
			validation.In("editor", "viewer").Error("Role must be either 'editor' or 'viewer'."),
		),
	)
}

type NoteMemberParams struct {
	ID       uuid.UUID `param:"id"`
	MemberID uuid.UUID `param:"memberID"`
}

func (p NoteMemberParams) New() interface{} {
	return &NoteMemberParams{}
}
