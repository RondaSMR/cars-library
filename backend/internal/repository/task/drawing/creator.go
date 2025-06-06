package drawing

import (
	"backend/internal/domain/entities"
	"context"
)

func (r Repository) CreateDrawing(ctx context.Context, task entities.Drawing) (entities.Drawing, error) {
	_, err := r.pool.Exec(ctx, `insert into drawings(id, title, file_url, car_model, category, user_id,  uploaded_by) values($1, $2, $3, $4, $5, $6, %7)`,
		task.DrawingID, task.Title, task.FileUrl, task.CarModel, task.Category, task.UserID, task.UploadedBy)
	if err != nil {
		return entities.Drawing{}, err
	}

	return task, nil
}
