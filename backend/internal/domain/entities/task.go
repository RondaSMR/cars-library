package entities

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	CommID    uuid.UUID
	Content   string
	CreatedAt time.Time
	UserID    uuid.UUID
	Username  string
}

type Drawing struct {
	DrawingID  uuid.UUID
	Title      string
	FileUrl    string
	CarModel   string
	Category   string
	UserID     uuid.UUID
	UploadedBy time.Time
	Comments   []Comment
}

type NewComment struct {
	CommID    uuid.UUID
	UserID    uuid.UUID
	BookID    uuid.UUID
	DrawingID uuid.UUID
	Username  string
	Content   string
	CreatedAt time.Time
}

type Book struct {
	BookID     uuid.UUID
	Title      string
	FileUrl    string
	CarModel   string
	Category   string
	UserID     uuid.UUID
	UploadedBy time.Time
	Comments   []Comment
}
