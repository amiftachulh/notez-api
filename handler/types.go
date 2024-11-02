package handler

import (
	"time"

	"github.com/google/uuid"
)

type user struct {
	ID        uuid.UUID `json:"id"`
	Name      *string   `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type note struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Content   *string `json:"content"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type noteInviter struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  *string   `json:"name"`
}

type noteInvitationNote struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type noteInvitationResponse struct {
	ID        uuid.UUID          `json:"id"`
	Note      noteInvitationNote `json:"note"`
	Inviter   noteInviter        `json:"inviter"`
	Role      string             `json:"role"`
	CreatedAt time.Time          `json:"created_at"`
}
