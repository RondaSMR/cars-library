package drawing

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

var _ website.DrawingUsecase = new(Usecase)

type Repository interface {
	CreateDrawing(ctx context.Context, drawing entities.Drawing) (entities.Drawing, error)
	GetDrawing(ctx context.Context, id uuid.UUID) (entities.Drawing, error)
	DeleteDrawing(ctx context.Context, id uuid.UUID) error
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

func (u Usecase) CreateDrawing(ctx context.Context, task *entities.Drawing, fileData []byte) (entities.Drawing, error) {
	task.UploadedBy = time.Now()

	if err := u.messageQueue.Publish(
		ctx,
		u.infoMessageQueue,
		[]byte(fmt.Sprintf(
			"Smbd create new drawing: %v", task.DrawingID)),
	); err != nil {
		return entities.Drawing{}, err
	}

	// Создание имени для S3 и базы данных
	fileID := uuid.New()
	task.DrawingID = fileID
	fileName := fileID.String()

	// Загрузка файла в S3
	fileURL, err := u.s3Client.UploadFile(ctx, fileName, fileData)
	if err != nil {
		return entities.Drawing{}, fmt.Errorf("failed to upload file to S3: %w", err)
	}
	task.FileUrl = fileURL

	returnedTask, err := u.taskRepository.CreateDrawing(ctx, *task)
	if err != nil {
		// Если ошибка при отправке в базу данных, удаляем файл из S3
		if s3err := u.s3Client.DeleteFile(ctx, fileName); s3err != nil {
			return entities.Drawing{}, fmt.Errorf("%w and failed to delete file from S3: %w", err, s3err)
		}

		return returnedTask, err
	}

	return returnedTask, nil
}

func (u Usecase) GetDrawing(ctx context.Context, id uuid.UUID) (entities.Drawing, error) {
	if err := u.messageQueue.Publish(
		ctx,
		u.infoMessageQueue,
		[]byte(fmt.Sprintf(
			"Smbd get drawing: %v", id)),
	); err != nil {
		return entities.Drawing{}, err
	}
	return u.taskRepository.GetDrawing(ctx, id)
}

func (u Usecase) DeleteDrawing(ctx context.Context, id uuid.UUID) error {
	if err := u.messageQueue.Publish(
		ctx,
		u.infoMessageQueue,
		[]byte(fmt.Sprintf(
			"Smbd deleted drawing: %v", id)),
	); err != nil {
		return err
	}

	// Удаляем файл из базы данных
	if err := u.taskRepository.DeleteDrawing(ctx, id); err != nil {
		return err
	}

	// Удаляем файл из S3
	fileName := id.String()
	if err := u.s3Client.DeleteFile(ctx, fileName); err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}
