package drawing

import (
	"backend/internal/apperor"
	"context"
	"github.com/google/uuid"
)

func (r *Repository) DeleteDrawing(ctx context.Context, id uuid.UUID) error {
	drawingToDel, err := r.pool.Exec(ctx, `delete from drawings where id = $1`, id)
	if err != nil {
		return err
	}

	if drawingToDel.RowsAffected() == 0 {
		return apperor.ErrNoEffect
	}

	return nil
}
