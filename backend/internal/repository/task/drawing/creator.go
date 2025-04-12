package drawing

import (
	"backend/internal/domain/entities"
	"context"
)

func (r Repository) CreateDrawing(ctx context.Context, task entities.Drawing) (entities.Drawing, error) {
	_, err := r.pool.Exec(ctx, `insert into drawings(id, title, file_url, car_model, category, uploaded_by) values($1, $2, $3, $4, $5, $6)`,
		task.DrawingID, task.Title, task.FileUrl, task.CarModel, task.Category, task.UploadedBy)
	if err != nil {
		return entities.Drawing{}, err
	}

	return task, nil
}
