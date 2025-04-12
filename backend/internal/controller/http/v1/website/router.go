package website

import (
	"backend/internal/domain/entities"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DrawingUsecase interface {
	CreateDrawing(ctx context.Context, task *entities.Drawing, fileData []byte) (entities.Drawing, error)
	GetDrawing(ctx context.Context, id uuid.UUID) (entities.Drawing, error)
	DeleteDrawing(ctx context.Context, id uuid.UUID) error
}

type CommentUsecase interface {
	CreateComment(ctx context.Context, task *entities.NewComment) (entities.NewComment, error)
	DeleteComment(ctx context.Context, id uuid.UUID) error
}

type BookUsecase interface {
	// TODO: implement me
}

type router struct {
	drawingUsecase DrawingUsecase
	commentUsecase CommentUsecase
	bookUsecase    BookUsecase
}

func Router(
	ginGroup *gin.RouterGroup,
	drawingUsecase DrawingUsecase,
	commentUsecase CommentUsecase,
	bookUsecase BookUsecase,
	user string,
	pass string,
) {
	drawingRouter := router{drawingUsecase: drawingUsecase}
	commentRouter := router{commentUsecase: commentUsecase}
	//bookRouter := router{bookUsecase: bookUsecase}

	ginGroup.Use(gin.BasicAuth(gin.Accounts{
		user: pass,
	}))

	// For Drawings
	ginGroup.POST("/drawing", drawingRouter.CreateDrawing)
	ginGroup.GET("/drawing", drawingRouter.GetDrawing)
	ginGroup.DELETE("/drawing", drawingRouter.DeleteDrawing)

	// For Comments
	ginGroup.POST("/comment", commentRouter.CreateComment)
	ginGroup.DELETE("/comment", commentRouter.DeleteComment)

	// For Books
	// TODO: implement me
}
