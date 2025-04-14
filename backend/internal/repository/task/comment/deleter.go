package comment

import (
	"backend/internal/apperor"
	"context"
	"github.com/google/uuid"
)

func (r *Repository) DeleteComment(ctx context.Context, id uuid.UUID) error {
	comment, err := r.pool.Exec(ctx, `delete from comments where id = $1`, id)
	if err != nil {
		return err
	}

	if comment.RowsAffected() == 0 {
		return apperor.ErrNoEffect
	}

	return nil
}
