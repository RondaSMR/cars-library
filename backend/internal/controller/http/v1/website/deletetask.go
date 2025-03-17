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

	if err := r.taskUsecase.DeleteDrawing(c.Request.Context(), drawingID); err != nil {
		apperor.ErrInternalSystem.JsonResponse(c, err)
		return
	}

	c.Status(http.StatusNoContent)
	return
}
