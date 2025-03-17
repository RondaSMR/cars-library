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
	UploadedBy time.Time
	Comments   []Comment
}
