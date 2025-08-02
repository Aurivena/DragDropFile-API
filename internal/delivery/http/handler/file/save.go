package file

import (
	"github.com/Aurivena/answer"
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
func (r *Route) SaveFile(c *gin.Context) {
	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		logrus.Error("missing session ID header")
		answer.SendResponseSuccess(c, nil, answer.BadRequest)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		logrus.WithError(err).Error("failed to parse multipart form")
		answer.SendResponseSuccess(c, nil, answer.BadRequest)
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

	output, processStatus := r.action.SaveFiles(context.Background(), sessionID, files, headers)
	answer.SendResponseSuccess(c, output, processStatus)
}
