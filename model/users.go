package model

import (
	"regexp"

	"github.com/google/uuid"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      *string   `json:"name"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"-"`
	Role      string    `json:"role,omitempty"`
	CreatedAt string    `json:"created_at,omitempty"`
	UpdatedAt string    `json:"updated_at,omitempty"`
}

type UpdateUserInfo struct {
	Name *string `json:"name,omitempty"`
}

func (u UpdateUserInfo) New() interface{} {
	return &UpdateUserInfo{}
}

func (u UpdateUserInfo) Validate() error {
	return validation.ValidateStruct(
		&u,
		validation.Field(
			&u.Name,
			validation.When(
				u.Name != nil,
				validation.RuneLength(1, 256).Error("Name must be between 1 and 256 characters."),
			),
		),
	)
}

type UpdateUserEmail struct {
	Email string `json:"email"`
}

func (u UpdateUserEmail) New() interface{} {
	return &UpdateUserEmail{}
}

func (u UpdateUserEmail) Validate() error {
	return validation.ValidateStruct(
		&u,
		validation.Field(
			&u.Email,
			validation.Required.Error("Email is required."),
			is.Email.Error("Email is not valid."),
		),
	)
}

type UpdateUserPassword struct {
	CurrentPassword string `json:"current_password"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (u UpdateUserPassword) New() interface{} {
	return &UpdateUserPassword{}
}

func (u UpdateUserPassword) Validate() error {
	return validation.ValidateStruct(
		&u,
		validation.Field(
			&u.CurrentPassword,
			validation.Required.Error("Current password is required."),
		),
		validation.Field(
			&u.Password,
			validation.Required.Error("Password is required."),
			validation.RuneLength(8, 64).Error("Password must be between 8 and 64 characters."),
			validation.Match(regexp.MustCompile(`^[^\p{Cc}]+$`)).
				Error("Password can't contain invalid characters."),
		),
		validation.Field(&u.ConfirmPassword, validation.By(samePassword(u.Password))),
	)
}
