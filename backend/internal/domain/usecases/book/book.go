package book

import (
	"backend/internal/controller/http/v1/website"
	"backend/internal/domain/entities"
	"backend/pkg/connectors/S3connector"
	"backend/pkg/connectors/rabbitconnector"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

var _ website.BookUsecase = new(Usecase)

type Repository interface {
	CreateBook(ctx context.Context, book entities.Book) (entities.Book, error)
	GetBook(ctx context.Context, id uuid.UUID) (entities.Book, error)
	DeleteBook(ctx context.Context, id uuid.UUID) error
}

type Usecase struct {
	taskRepository   Repository
	infoMessageQueue string
	messageQueue     rabbitconnector.Connector
	s3Client         S3connector.Client
	s3Bucket         string
}

func NewUseCase(
	taskRepository Repository,
	mqConnector rabbitconnector.Connector,
	s3Client S3connector.Client,
	s3Bucket string,
	infoMessageQueue string,
) *Usecase {
	return &Usecase{
		infoMessageQueue: infoMessageQueue,
		taskRepository:   taskRepository,
		messageQueue:     mqConnector,
		s3Client:         s3Client,
		s3Bucket:         s3Bucket,
	}
}

func (u Usecase) CreateBook(ctx context.Context, task *entities.Book, fileData []byte) (entities.Book, error) {
	if err := u.messageQueue.Publish(
		ctx,
		u.infoMessageQueue,
		[]byte(fmt.Sprintf(
			"%v create new book", task.UserID)),
	); err != nil {
		return entities.Book{}, err
	}

	task.UploadedBy = time.Now()

	// Создание имени для S3 и базы данных
	fileID := uuid.New()
	task.BookID = fileID
	fileName := fileID.String()

	// Загрузка файла в S3
	fileURL, err := u.s3Client.UploadFile(ctx, fileName, fileData)
	if err != nil {
		return entities.Book{}, fmt.Errorf("failed to upload book to S3: %w", err)
	}
	task.FileUrl = fileURL

	returnedTask, err := u.taskRepository.CreateBook(ctx, *task)
	if err != nil {
		// Если ошибка при отправке в базу данных, удаляем файл из S3
		if s3err := u.s3Client.DeleteFile(ctx, fileName); s3err != nil {
			return entities.Book{}, fmt.Errorf("%w and failed to delete book from S3: %w", err, s3err)
		}

		return returnedTask, err
	}

	return returnedTask, nil
}

func (u Usecase) GetBook(ctx context.Context, id uuid.UUID) (entities.Book, error) {
	if err := u.messageQueue.Publish(
		ctx,
		u.infoMessageQueue,
		[]byte(fmt.Sprintf(
			"Smbd get book: %v", id)),
	); err != nil {
		return entities.Book{}, err
	}
	return u.taskRepository.GetBook(ctx, id)
}

func (u Usecase) DeleteBook(ctx context.Context, id uuid.UUID) error {
	if err := u.messageQueue.Publish(
		ctx,
		u.infoMessageQueue,
		[]byte(fmt.Sprintf(
			"Smbd deleted book: %v", id)),
	); err != nil {
		return err
	}

	// Удаляем файл из базы данных
	if err := u.taskRepository.DeleteBook(ctx, id); err != nil {
		return err
	}

	// Удаляем файл из S3
	fileName := id.String()
	if err := u.s3Client.DeleteFile(ctx, fileName); err != nil {
		return fmt.Errorf("failed to delete book from S3: %w", err)
	}

	return nil
}
