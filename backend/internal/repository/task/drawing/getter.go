package drawing

import (
	TaskDAO "backend/internal/adapters/repository"
	"backend/internal/apperor"
	"backend/internal/domain/entities"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r Repository) GetDrawing(ctx context.Context, id uuid.UUID) (entities.Drawing, error) {
	drawingRow, err := r.pool.Query(ctx, `
		select id, title, file_url, car_model, category, uploaded_by
		from drawings
		where drawing_id = $1
	`, id)
	if err != nil {
		return entities.Drawing{}, fmt.Errorf("executing query error: %w", err)
	}

	drawingRepo, err := pgx.CollectOneRow(drawingRow, pgx.RowToStructByName[TaskDAO.Drawing])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Drawing{}, apperor.ErrRepoNotFound
		}
		return entities.Drawing{}, fmt.Errorf("selecting drawing from database error: %w", err)
	}

	commentRows, err := r.pool.Query(ctx, `
		select id, content, created_at, user_id, username
		from comments
		where drawing_id = $1
	`, id)
	if err != nil {
		return entities.Drawing{}, fmt.Errorf("executing comment query error: %w", err)
	}

	commentsRepo, err := pgx.CollectRows(commentRows, pgx.RowToStructByName[TaskDAO.Comment])
	if err != nil {
		return entities.Drawing{}, fmt.Errorf("selecting comments from database error: %w", err)
	}

	drawing := TaskDAO.AdapterRepoDrawingToEntity(drawingRepo)
	for _, commentRepo := range commentsRepo {
		drawing.Comments = append(drawing.Comments, TaskDAO.AdapterRepoCommentToEntity(commentRepo))
	}

	return drawing, nil
}
