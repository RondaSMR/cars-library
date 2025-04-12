package website

import (
	"backend/internal/adapters/controller/http/taskapi"
	"backend/internal/apperor"
	"backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io"
	"net/http"
)

func (r router) CreateDrawing(c *gin.Context) {
	// Получаем JSON как строку из формы
	jsonData := c.PostForm("json")

	// Декодируем JSON в структуру
	var taskInput taskapi.Drawing
	if err := json.Unmarshal([]byte(jsonData), &taskInput); err != nil {
		apperor.ErrBadRequest.JsonResponse(c, err)
		return
	}

	// Получаем файл (две переменные: header и error)
	header, err := c.FormFile("drawing")
	if err != nil {
		apperor.ErrBadRequest.JsonResponse(c, err)
		return
	}

	// Открываем сам файл
	src, err := header.Open()
	if err != nil {
		apperor.ErrInternalSystem.JsonResponse(c, err)
		return
	}
	defer src.Close()

	// Читаем в []byte
	fileData, err := io.ReadAll(src)
	if err != nil {
		apperor.ErrInternalSystem.JsonResponse(c, err)
		return
	}

	// Конвертируем в сущность
	pointerTask := taskapi.AdapterHttpDrawingToEntity(taskInput)

	// Передаём структуру и файл в usecase
	drawing, err := r.drawingUsecase.CreateDrawing(c.Request.Context(), &pointerTask, fileData)
	if err != nil {
		apperor.ErrInternalSystem.JsonResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, utils.GenerateResponse(nil, taskapi.AdapterEntityToHttpDrawing(drawing)))
}

func (r router) CreateComment(c *gin.Context) {
	var taskInput taskapi.NewComment
	if err := c.ShouldBindJSON(&taskInput); err != nil {
		apperor.ErrBadRequest.JsonResponse(c, err)
		return
	}

	pointerTask := taskapi.AdapterNewHttpCommentToEntity(taskInput)

	comment, err := r.commentUsecase.CreateComment(c.Request.Context(), &pointerTask)
	if err != nil {
		apperor.ErrInternalSystem.JsonResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, utils.GenerateResponse(nil, taskapi.AdapterNewEntityToHttpComment(comment)))
}
