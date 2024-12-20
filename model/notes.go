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

func (n NoteInput) New() interface{} {
	return &NoteInput{}
}

func (n NoteInput) Validate() error {
	return validation.ValidateStruct(
		&n,
		validation.Field(
			&n.Title,
			validation.Required.Error("Title is required."),
			validation.RuneLength(1, 300).Error("Title must be less than 300 characters."),
		),
		validation.Field(
			&n.Content,
			validation.When(
				n.Content != nil,
				validation.Required.Error("Content must be at least 1 character."),
				validation.Length(1, NOTE_MAX_CONTENT_BYTES).
					Error("Content size must be less than 5 MB."),
			),
		),
	)
}

type NoteParams struct {
	ID uuid.UUID `param:"id"`
}

func (p NoteParams) New() interface{} {
	return &NoteParams{}
}

type NoteQuery struct {
	Query    string `query:"q"         json:"q"`
	Page     int    `query:"page"      json:"page"`
	PageSize int    `query:"page_size" json:"page_size"`
	Sort     string `query:"sort"      json:"sort"`
	Order    string `query:"order"     json:"order"`
	Role     string `query:"role" json:"role"`
}

func (q NoteQuery) New() interface{} {
	return &NoteQuery{
		Page:     1,
		PageSize: 10,
		Sort:     "id",
		Order:    "asc",
	}
}

func (q NoteQuery) Validate() error {
	return validation.ValidateStruct(
		&q,
		validation.Field(
			&q.Page,
			validation.Min(1).Error("Page must be greater than 0."),
		),
		validation.Field(
			&q.PageSize,
			validation.Min(1).Error("Page size must be greater than 0."),
			validation.Max(100).Error("Page size must be less than 100."),
		),
		validation.Field(
			&q.Sort,
			validation.In("id", "title", "created_at", "updated_at").
				Error("Invalid sort field. Allowed fields: 'title', 'created_at', 'updated_at'."),
		),
		validation.Field(
			&q.Order,
			validation.In("asc", "desc").Error("Invalid order. Allowed values: 'asc', 'desc'."),
		),
		validation.Field(
			&q.Role,
			validation.In("owner", "editor", "viewer").
				Error("Invalid role. Allowed values: 'owner', 'editor', 'viewer'."),
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

type NoteMember struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      *string   `json:"name"`
	Role      string    `json:"role,omitempty"`
	CreatedAt string    `json:"created_at,omitempty"`
}

type NoteDetail struct {
	ID        uuid.UUID    `json:"id"`
	Title     string       `json:"title"`
	Content   *string      `json:"content"`
	Role      *string      `json:"role"`
	Owner     NoteMember   `json:"owner"`
	Members   []NoteMember `json:"members"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
}
