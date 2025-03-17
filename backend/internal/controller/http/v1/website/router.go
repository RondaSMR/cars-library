package website

import (
	"backend/internal/domain/entities"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Usecase interface {
	GetDrawing(ctx context.Context, id uuid.UUID) (entities.Drawing, error)
	DeleteDrawing(ctx context.Context, id uuid.UUID) error
}

type router struct {
	taskUsecase Usecase
}

func Router(
	ginGroup *gin.RouterGroup,
	taskUsecase Usecase,
	user string,
	pass string,
) {
	r := router{taskUsecase: taskUsecase}

	ginGroup.Use(gin.BasicAuth(gin.Accounts{
		user: pass,
	}))

	// TODO: ginGroup.(POST, GET, DELETE)
	ginGroup.GET("", r.GetDrawing)
}
