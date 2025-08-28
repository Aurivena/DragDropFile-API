package file

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"mime/multipart"
)

// @Summary      Сохранить файл
// @Description  Сохраняет файл и параметры, переданные пользователем.
// @Tags         File
// @Accept       multipart/form-data
// @Produce      json
// @Param        X-Session-ID header string true "Идентификатор сессии пользователя"
// @Param        files formData file true "Файл для загрузки"
// @Success      200 {object} entity.FileSaveOutput "Файл успешно сохранен"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/save [post]
func (h *Handler) Execute(c *gin.Context) {
	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		logrus.Error("missing session ID header")
		SendResponseSuccess(c, nil, BadRequest)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		logrus.WithError(err).Error("failed to parse multipart form")
		SendResponseSuccess(c, nil, BadRequest)
		return
	}

	headers := form.File["files"]
	var files []multipart.File

	for _, header := range headers {
		f, err := header.Open()
		if err != nil {
			logrus.WithError(err).Error("failed to open uploaded file")
			continue
		}
		files = append(files, f)
		defer f.Close()
	}

	output, processStatus := h.application.Execute(context.Background(), sessionID, files, headers)
	SendResponseSuccess(c, output, processStatus)
}
