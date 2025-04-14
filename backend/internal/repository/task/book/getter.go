package book

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

func (r Repository) GetBook(ctx context.Context, id uuid.UUID) (entities.Book, error) {
	bookRow, err := r.pool.Query(ctx, `
		select id, title, file_url, car_model, category, uploaded_by
		from books
		where book_id = $1
	`, id)
	if err != nil {
		return entities.Book{}, fmt.Errorf("executing query error: %w", err)
	}

	bookRepo, err := pgx.CollectOneRow(bookRow, pgx.RowToStructByName[TaskDAO.Book])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Book{}, apperor.ErrRepoNotFound
		}
		return entities.Book{}, fmt.Errorf("selecting book from database error: %w", err)
	}

	commentRows, err := r.pool.Query(ctx, `
		select id, content, created_at, user_id, username
		from comments
		where book_id = $1
	`, id)
	if err != nil {
		return entities.Book{}, fmt.Errorf("executing comment query error: %w", err)
	}

	commentsRepo, err := pgx.CollectRows(commentRows, pgx.RowToStructByName[TaskDAO.Comment])
	if err != nil {
		return entities.Book{}, fmt.Errorf("selecting comments from database error: %w", err)
	}

	book := TaskDAO.AdapterRepoBookToEntity(bookRepo)
	for _, commentRepo := range commentsRepo {
		book.Comments = append(book.Comments, TaskDAO.AdapterRepoCommentToEntity(commentRepo))
	}

	return book, nil
}
