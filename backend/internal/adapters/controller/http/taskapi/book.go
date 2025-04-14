package taskapi

import (
	"backend/internal/domain/entities"
	"github.com/google/uuid"
	"time"
)

type Book struct {
	BookID     uuid.UUID `json:"book_id"`
	Title      string    `json:"title"`
	FileUrl    string    `json:"file_url"`
	CarModel   string    `json:"car_model"`
	Category   string    `json:"category"`
	UserID     uuid.UUID `json:"user_id"`
	UploadedBy time.Time `json:"uploaded_by,omitempty"`
	Comments   []Comment `json:"comments,omitempty"`
}

func AdapterEntityToHttpBook(book entities.Book) Book {
	comments := make([]Comment, len(book.Comments))
	for i, comment := range book.Comments {
		comments[i] = Comment{
			CommID:    comment.CommID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
			UserID:    comment.UserID,
			Username:  comment.Username,
		}
	}

	return Book{
		BookID:     book.BookID,
		Title:      book.Title,
		FileUrl:    book.FileUrl,
		CarModel:   book.CarModel,
		Category:   book.Category,
		UserID:     book.UserID,
		UploadedBy: book.UploadedBy,
		Comments:   comments,
	}
}

func AdapterHttpBookToEntity(book Book) entities.Book {
	return entities.Book{
		BookID:     book.BookID,
		Title:      book.Title,
		FileUrl:    book.FileUrl,
		CarModel:   book.CarModel,
		Category:   book.Category,
		UserID:     book.UserID,
		UploadedBy: book.UploadedBy,
	}
}
