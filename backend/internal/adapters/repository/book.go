package repository

import (
	"backend/internal/domain/entities"
	"database/sql"
	"github.com/google/uuid"
)

type Book struct {
	BookID     uuid.NullUUID  `db:"book_id"`
	Title      sql.NullString `db:"title"`
	FileUrl    sql.NullString `db:"file_url"`
	CarModel   sql.NullString `db:"car_model"`
	Category   sql.NullString `db:"category"`
	UserID     uuid.NullUUID  `db:"user_id"`
	UploadedBy sql.NullTime   `db:"uploaded_by"`
	Comments   []Comment
}

func AdapterRepoBookToEntity(book Book) entities.Book {
	return entities.Book{
		BookID:     book.BookID.UUID,
		Title:      book.Title.String,
		FileUrl:    book.FileUrl.String,
		CarModel:   book.CarModel.String,
		Category:   book.Category.String,
		UserID:     book.UserID.UUID,
		UploadedBy: book.UploadedBy.Time,
	}
}
