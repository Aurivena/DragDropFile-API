package file

import (
	"DragDrop-Files/internal/middleware"
	"mime/multipart"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Tags         File
// @Summary      Загрузить файл(ы)
// @Description  Принимает один или несколько файлов в multipart/form-data и сохраняет их.
// @Accept       multipart/form-data
// @Produce      json
// @Param        X-Session-ID  header   string  true  "Идентификатор сессии пользователя"
// @Param        files         formData file    true  "Файлы для загрузки" collectionFormat(multi)
// @Success      200           {object} entity.FileSaveOutput  "Файл(ы) успешно сохранён(ы)"
// @Failure      400           {object} map[string]any         "Некорректные данные (Spond error)"
// @Failure      500           {object} map[string]any         "Внутренняя ошибка сервера (Spond error)"
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

	if len(files) == 0 || len(files) != len(headers) {
		h.spond.SendResponseError(c.Writer, h.BadRequest("1. Ваша сессия недействительна\n"+"2. Длина загруженных файлов == 0"))
		return
	}

	output, errResp := h.application.File.Execute(c.GetHeader(middleware.Session), files, headers)
	if errResp != nil {
		h.spond.SendResponseError(c.Writer, errResp)
		return
	}
	h.spond.SendResponseSuccess(c.Writer, envelope.Success, output)
}
