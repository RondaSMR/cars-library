package book

import (
	"backend/internal/domain/entities"
	"context"
)

func (r Repository) CreateBook(ctx context.Context, task entities.Book) (entities.Book, error) {
	_, err := r.pool.Exec(ctx, `insert into book(id, title, file_url, car_model, category, user_id, uploaded_by) values($1, $2, $3, $4, $5, $6, $7)`,
		task.BookID, task.Title, task.FileUrl, task.CarModel, task.Category, task.UserID, task.UploadedBy)
	if err != nil {
		return entities.Book{}, err
	}

	return task, nil
}
