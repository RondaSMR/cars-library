package book

import (
	"backend/internal/controller/http/v1/website"
	"backend/pkg/connectors/rabbitconnector"
)

var _ website.BookUsecase = new(Usecase)

type Repository interface {
	// TODO: implement me
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
