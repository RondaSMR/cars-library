package repository

import (
	"backend/internal/domain/entities"
	"database/sql"
	"github.com/google/uuid"
)

type Drawing struct {
	DrawingID  uuid.NullUUID  `db:"drawing_id"`
	Title      sql.NullString `db:"title"`
	FileUrl    sql.NullString `db:"file_url"`
	CarModel   sql.NullString `db:"car_model"`
	Category   sql.NullString `db:"category"`
	UserID     uuid.NullUUID  `db:"user_id"`
	UploadedBy sql.NullTime   `db:"uploaded_by"`
	Comments   []Comment
}

func AdapterRepoDrawingToEntity(drawing Drawing) entities.Drawing {
	return entities.Drawing{
		DrawingID:  drawing.DrawingID.UUID,
		Title:      drawing.Title.String,
		FileUrl:    drawing.FileUrl.String,
		CarModel:   drawing.CarModel.String,
		Category:   drawing.Category.String,
		UserID:     drawing.UserID.UUID,
		UploadedBy: drawing.UploadedBy.Time,
	}
}
