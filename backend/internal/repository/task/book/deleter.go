package book

import (
	"backend/internal/apperor"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"log"
)

func (r *Repository) DeleteBook(ctx context.Context, id uuid.UUID) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("rollback error: %v", err)
		}
	}()

	// Удаление комментариев
	if _, err = tx.Exec(ctx, `delete from comments where book_id=$1`, id); err != nil {
		return err
	}

	// Удаление рисунка
	tag, err := tx.Exec(ctx, `delete from books where id=$1`, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return apperor.ErrNoEffect
	}

	return tx.Commit(ctx)
}
