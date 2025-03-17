package taskapi

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	CommID    uuid.UUID `json:"comm_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
}
