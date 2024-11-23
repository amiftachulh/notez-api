package model

import "github.com/invopop/validation"

type UpdateNoteMemberRole struct {
	Role string `json:"role"`
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
