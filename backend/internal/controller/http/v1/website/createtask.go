package website

import (
	"backend/internal/adapters/controller/http/taskapi"
	"backend/internal/apperor"
	"backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (r router) CreateDrawing(c *gin.Context) {
	var taskInput taskapi.Drawing

	if err := c.ShouldBindJSON(&taskInput); err != nil {
		apperor.ErrBadRequest.JsonResponse(c, err)
		return
	}

	pointerTask := taskapi.AdapterHttpDrawingToEntity(taskInput)

	drawing, err := r.taskUsecase.CreateDrawing(c.Request.Context(), &pointerTask)
	if err != nil {
		apperor.ErrInternalSystem.JsonResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, utils.GenerateResponse(nil, taskapi.AdapterEntityToHttpDrawing(drawing)))
}
