package website

import (
	"backend/internal/adapters/controller/http/taskapi"
	"backend/internal/apperor"
	"backend/pkg/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (r router) GetDrawing(c *gin.Context) {
	drawingID, err := uuid.Parse(c.Param("drawing_id"))
	if err != nil {
		apperor.ErrInvalidID.JsonResponse(c, err)
		return
	}

	drawing, err := r.drawingUsecase.GetDrawing(c.Request.Context(), drawingID)
	if err != nil {
		if errors.Is(err, apperor.ErrRepoNotFound) {
			apperor.ErrNotFound.JsonResponse(c, err)
			return
		}
		apperor.ErrInternalSystem.JsonResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.GenerateResponse(nil, taskapi.AdapterEntityToHttpDrawing(drawing)))
}
