package task

import (
	"backend/internal/controller/http/v1/website"
	"backend/internal/domain/entities"
	"backend/pkg/connectors/rabbitconnector"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

var _ website.Usecase = new(Usecase)

type Repository interface {
	CreateDrawing(ctx context.Context, drawing entities.Drawing) (entities.Drawing, error)
	GetDrawing(ctx context.Context, id uuid.UUID) (entities.Drawing, error)
	DeleteDrawing(ctx context.Context, id uuid.UUID) error

	// TODO : ReturnWebsites...
}

type Usecase struct {
	taskRepository   Repository
	infoMessageQueue string
	messageQueue     rabbitconnector.Connector
}

func NewUseCase(
	taskRepository Repository,
	mqConnector rabbitconnector.Connector,
	infoMessageQueue string,
) *Usecase {
	return &Usecase{
		infoMessageQueue: infoMessageQueue,
		taskRepository:   taskRepository,
		messageQueue:     mqConnector,
	}
}

// / FOR DRAWINGS
func (u Usecase) CreateDrawing(ctx context.Context, task *entities.Drawing) (entities.Drawing, error) {
	task.DrawingID = uuid.New()
	task.UploadedBy = time.Now()

	if err := u.messageQueue.Publish(
		ctx,
		u.infoMessageQueue,
		[]byte(fmt.Sprintf(
			"Smbd create new drawing: %v", task.DrawingID)),
	); err != nil {
		return entities.Drawing{}, err
	}
	return u.taskRepository.CreateDrawing(ctx, *task)
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
	return u.taskRepository.DeleteDrawing(ctx, id)
}
