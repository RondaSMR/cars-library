package taskapi

import (
	"backend/internal/domain/entities"
	"github.com/google/uuid"
	"time"
)

type Drawing struct {
	DrawingID  uuid.UUID `json:"drawing_id,omitempty"`
	Title      string    `json:"title"`
	FileUrl    string    `json:"file_url"`
	CarModel   string    `json:"car_model"`
	Category   string    `json:"category"`
	UploadedBy time.Time `json:"uploaded_by,omitempty"`
	Comments   []Comment `json:"comments,omitempty"`
}

func AdapterEntityToHttpDrawing(drawing entities.Drawing) Drawing {
	comments := make([]Comment, len(drawing.Comments))
	for i, comment := range drawing.Comments {
		comments[i] = Comment{
			CommID:    comment.CommID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
			UserID:    comment.UserID,
			Username:  comment.Username,
		}
	}

	return Drawing{
		DrawingID:  drawing.DrawingID,
		Title:      drawing.Title,
		FileUrl:    drawing.FileUrl,
		CarModel:   drawing.CarModel,
		Category:   drawing.Category,
		UploadedBy: drawing.UploadedBy,
		Comments:   comments,
	}
}

func AdapterHttpDrawingToEntity(drawing Drawing) entities.Drawing {
	return entities.Drawing{
		DrawingID:  drawing.DrawingID,
		Title:      drawing.Title,
		FileUrl:    drawing.FileUrl,
		CarModel:   drawing.CarModel,
		Category:   drawing.Category,
		UploadedBy: drawing.UploadedBy,
	}
}
