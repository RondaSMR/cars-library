package taskapi

import (
	"backend/internal/domain/entities"
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

type NewComment struct {
	CommID    uuid.UUID `json:"comm_id,omitempty"`
	UserID    uuid.UUID `json:"user_id"`
	BookID    uuid.UUID `json:"book_id,omitempty"`
	DrawingID uuid.UUID `json:"drawing_id,omitempty"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func AdapterNewHttpCommentToEntity(comment NewComment) entities.NewComment {
	return entities.NewComment{
		CommID:    comment.CommID,
		UserID:    comment.UserID,
		BookID:    comment.BookID,
		DrawingID: comment.DrawingID,
		Username:  comment.Username,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}
}

func AdapterNewEntityToHttpComment(comment entities.NewComment) NewComment {
	return NewComment{
		CommID:    comment.CommID,
		UserID:    comment.UserID,
		BookID:    comment.BookID,
		DrawingID: comment.DrawingID,
		Username:  comment.Username,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}
}
