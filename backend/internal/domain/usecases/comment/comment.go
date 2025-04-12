package comment

import (
	"backend/internal/controller/http/v1/website"
	"backend/internal/domain/entities"
	"backend/pkg/connectors/rabbitconnector"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

var _ website.CommentUsecase = new(Usecase)

type Repository interface {
	CreateComment(ctx context.Context, comment entities.NewComment) (entities.NewComment, error)
	DeleteComment(ctx context.Context, id uuid.UUID) error
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

func (u Usecase) CreateComment(ctx context.Context, task *entities.NewComment) (entities.NewComment, error) {
	task.CommID = uuid.New()
	task.CreatedAt = time.Now()

	if err := u.messageQueue.Publish(
		ctx,
		u.infoMessageQueue,
		[]byte(fmt.Sprintf(
			"Smbd create new comment: %v", task.CommID)),
	); err != nil {
		return entities.NewComment{}, err
	}
	return u.taskRepository.CreateComment(ctx, *task)
}

func (u Usecase) DeleteComment(ctx context.Context, id uuid.UUID) error {
	if err := u.messageQueue.Publish(
		ctx,
		u.infoMessageQueue,
		[]byte(fmt.Sprintf(
			"Smbd deleted comment: %v", id)),
	); err != nil {
		return err
	}
	return u.taskRepository.DeleteComment(ctx, id)
}
