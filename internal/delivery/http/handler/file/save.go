package file

import (
	"DragDrop-Files/internal/middleware"
	"mime/multipart"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary      Сохранить файл
// @Description  Сохраняет файл и параметры, переданные пользователем.
// @Tags         Get
// @Accept       multipart/form-data
// @Produce      json
// @Param        X-SessionID-ID header string true "Идентификатор сессии пользователя"
// @Param        files formData file true "Файл для загрузки"
// @Success      200 {object} entity.FileSaveOutput "Файл успешно сохранен"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/save [post]
func (h *Handler) Execute(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		logrus.WithError(err).Error("failed to parse multipart form")
		h.spond.SendResponseError(c.Writer, h.ErrorSystem())
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

	output, errResp := h.application.File.Execute(c.GetHeader(middleware.Session), files, headers)
	if errResp != nil {
		h.spond.SendResponseError(c.Writer, errResp)
		return
	}
	h.spond.SendResponseSuccess(c.Writer, envelope.Success, output)
}
