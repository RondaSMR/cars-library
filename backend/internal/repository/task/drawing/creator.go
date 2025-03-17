package drawing

import (
	"backend/internal/domain/entities"
	"context"
)

func (r Repository) CreateDrawing(ctx context.Context, task entities.Drawing) (entities.Drawing, error) {
	_, err := r.pool.Exec(ctx, `insert into drawings(title, file_url, car_model, category) values($1, $2, $3, $4)`, task.Title, task.FileUrl, task.Category)
	if err != nil {
		return entities.Drawing{}, err
	}

	return task, nil
}
