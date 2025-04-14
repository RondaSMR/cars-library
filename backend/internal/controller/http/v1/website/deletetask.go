package website

import (
	"backend/internal/apperor"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (r router) DeleteDrawing(c *gin.Context) {
	drawingID, err := uuid.Parse(c.Query("drawing_id"))
	if err != nil {
		apperor.ErrInvalidID.JsonResponse(c, err)
		return
	}

	if err := r.drawingUsecase.DeleteDrawing(c.Request.Context(), drawingID); err != nil {
		apperor.ErrInternalSystem.JsonResponse(c, err)
		return
	}

	c.Status(http.StatusNoContent)
	return
}

// -----------------------------------------------

func (r router) DeleteComment(c *gin.Context) {
	commentID, err := uuid.Parse(c.Query("comment_id"))
	if err != nil {
		apperor.ErrInvalidID.JsonResponse(c, err)
		return
	}

	if err := r.commentUsecase.DeleteComment(c.Request.Context(), commentID); err != nil {
		apperor.ErrInternalSystem.JsonResponse(c, err)
		return
	}

	c.Status(http.StatusNoContent)
	return
}

// -----------------------------------------------

func (r router) DeleteBook(c *gin.Context) {
	bookID, err := uuid.Parse(c.Query("book_id"))
	if err != nil {
		apperor.ErrInvalidID.JsonResponse(c, err)
		return
	}

	if err := r.bookUsecase.DeleteBook(c.Request.Context(), bookID); err != nil {
		apperor.ErrInternalSystem.JsonResponse(c, err)
		return
	}

	c.Status(http.StatusNoContent)
	return
}
