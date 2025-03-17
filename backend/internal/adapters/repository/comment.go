package repository

import (
	"backend/internal/domain/entities"
	"database/sql"
	"github.com/google/uuid"
)

type Comment struct {
	CommID    uuid.NullUUID  `db:"comm_id"`
	Content   sql.NullString `db:"content"`
	CreatedAt sql.NullTime   `db:"created_at"`
	UserID    uuid.NullUUID  `db:"user_id"`
	Username  sql.NullString `db:"username"`
}

func AdapterRepoCommentToEntity(comment Comment) entities.Comment {
	return entities.Comment{
		CommID:    comment.CommID.UUID,
		Content:   comment.Content.String,
		UserID:    comment.UserID.UUID,
		Username:  comment.Username.String,
		CreatedAt: comment.CreatedAt.Time,
	}
}
