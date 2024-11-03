package model

import (
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

type Register struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func samePassword(str string) validation.RuleFunc {
	return func(value interface{}) error {
		s, _ := value.(string)
		if s != str {
			return validation.NewError("passwords_not_match", "Passwords do not match.")
		}
		return nil
	}
}

func (r Register) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(
			&r.Email,
			validation.Required.Error("Email is required."),
			is.Email.Error("Email is invalid."),
		),
		validation.Field(
			&r.Password,
			validation.Required.Error("Password is required."),
			validation.RuneLength(8, 64).Error("Password must be between 8 and 64 characters."),
			validation.Match(regexp.MustCompile(`^[^\p{Cc}]+$`)).
				Error("Password can't contain invalid characters."),
		),
		validation.Field(&r.ConfirmPassword, validation.By(samePassword(r.Password))),
	)
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l Login) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Email, validation.Required.Error("Email is required.")),
		validation.Field(&l.Password, validation.Required.Error("Password is required.")),
	)
}

type AuthUser struct {
	ID        uuid.UUID `json:"id"`
	Name      *string   `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
