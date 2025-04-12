package comment

import (
	"backend/internal/domain/entities"
	"context"
)

func (r Repository) CreateComment(ctx context.Context, task entities.NewComment) (entities.NewComment, error) {
	_, err := r.pool.Exec(ctx, `
		insert into comments(id, user_id, book_id, drawing_id, username, content, created_at)
		values ($1, $2, $3, $4, $5, $6, $7)
		`, task.CommID, task.UserID, task.BookID, task.DrawingID, task.Username, task.Content, task.CreatedAt)
	if err != nil {
		return entities.NewComment{}, err
	}

	return task, nil
}
