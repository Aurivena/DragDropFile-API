package route

import (
	"DragDrop-Files/models"
	"context"
	"fmt"
	"github.com/Aurivena/answer"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"mime/multipart"
)

// @Summary      Сохранить файл
// @Description  Сохраняет файл и параметры, переданные пользователем.
// @Tags         File
// @Accept       multipart/form-data
// @Produce      json
// @Param        X-Session-ID header string true "Идентификатор сессии пользователя"
// @Param        file formData file true "Файл для загрузки"
// @Success      200 {object} models.FilSaveOutput "Файл успешно сохранен"
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

	headers := form.File["file"]
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

	output, processStatus := r.action.Create(context.Background(), sessionID, files, headers)
	answer.SendResponseSuccess(c, output, processStatus)
}

// @Summary      Получить данные файла
// @Description  Получает данные файла по указаному id.
// @Tags         File
// @Produce      json
// @Param        id path string true "Идентификатор файла"
// @Success      200 {object} models.DataOutput "Файл успешно сохранен"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/:id/data [get]
func (r *Route) GetDataFile(c *gin.Context) {
	id := c.Param("id")

	output, processStatus := r.action.GetDataFile(id)
	answer.SendResponseSuccess(c, output, processStatus)
}

// @Summary      Получить файл
// @Description  Получает файл по переданному идентификатору.
// @Tags         File
// @Accept       json
// @Produce      octet-stream
// @Param        id path string true "Идентификатор файла"
// @Param        X-Password header string true "Пароль для файлов"
// @Success      200 {file} string "Файл успешно получен"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      401 {object} string "Неверный пароль"
// @Failure      410 {object} string "Хранение файла закончено. Файл удален"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/:id [get]
func (r *Route) Get(c *gin.Context) {
	id := c.Param("id")
	password := c.GetHeader("X-Password")

	out, processStatus := r.action.GetFile(id, password)
	if processStatus != answer.OK {
		answer.SendResponseSuccess(c, nil, processStatus)
		return
	}

	objInfo, err := out.File.Stat()
	if err != nil {
		log.Printf("Ошибка Stat() для объекта %s: %v", id, err)
		answer.SendResponseSuccess(c, nil, answer.InternalServerError)
		out.File.Close()
		return
	}

	contentDisposition := fmt.Sprintf("attachment; filename=\"%s\"", out.Name)

	c.DataFromReader(answer.OK,
		objInfo.Size,
		objInfo.ContentType,
		out.File,
		map[string]string{
			"Content-Disposition": contentDisposition,
		},
	)
}

// @Summary      Обновить дату удаления файла
// @Description  Обновляет количество дней до автоматического удаления файла.
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        input body models.DayDeletedUpdateInput true "Данные для ввода"
// @Param        X-Session-ID header string true "Идентификатор сессии пользователя"
// @Success      204 {object} string "NoContent" "Выходные данные"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/update/deleted [put]
func (r *Route) UpdateCountDayToDeleted(c *gin.Context) {
	var input *models.DayDeletedUpdateInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		answer.SendResponseSuccess(c, nil, answer.BadRequest)
		return
	}

	sessionID := c.GetHeader("X-Session-ID")

	processStatus := r.action.UpdateDateDeleted(input.CountDayToDeleted, sessionID)
	answer.SendResponseSuccess(c, nil, processStatus)
}

// @Summary      Обновить пароль для файла
// @Description  Обновляет пароль, необходимый для доступа к файлу.
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        input body models.PasswordUpdateInput true "Данные для ввода"
// @Param        X-Session-ID header string true "Идентификатор сессии пользователя"
// @Success      204 {object} string "NoContent" "Выходные данные"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/update/password [put]
func (r *Route) UpdatePassword(c *gin.Context) {
	var input *models.PasswordUpdateInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		answer.SendResponseSuccess(c, nil, answer.BadRequest)
		return
	}

	sessionID := c.GetHeader("X-Session-ID")

	processStatus := r.action.UpdatePassword(input.Password, sessionID)
	answer.SendResponseSuccess(c, nil, processStatus)
}

// @Summary      Обновить количество загрузок файла
// @Description  Устанавливает новое ограничение по количеству скачиваний файла.
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        input body models.CountDownloadUpdateInput true "Данные для ввода"
// @Param        X-Session-ID header string true "Идентификатор сессии пользователя"
// @Success      204 {object} string "NoContent" "Выходные данные"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/update/count-download [put]
func (r *Route) UpdateCountDownload(c *gin.Context) {
	var input *models.CountDownloadUpdateInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		answer.SendResponseSuccess(c, nil, answer.BadRequest)
		return
	}

	sessionID := c.GetHeader("X-Session-ID")

	processStatus := r.action.UpdateCountDownload(input.CountDownload, sessionID)
	answer.SendResponseSuccess(c, nil, processStatus)
}

// @Summary      Обновить описание файла
// @Description  Устанавливает новое описание для файла.
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        input body models.DescriptionUpdateInput true "Данные для ввода"
// @Param        X-Session-ID header string true "Идентификатор сессии пользователя"
// @Success      204 {object} string "NoContent" "Выходные данные"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/update/description [put]
func (r *Route) UpdateDescription(c *gin.Context) {
	var input *models.DescriptionUpdateInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		answer.SendResponseSuccess(c, nil, answer.BadRequest)
		return
	}

	sessionID := c.GetHeader("X-Session-ID")

	processStatus := r.action.UpdateDescription(input.Description, sessionID)
	answer.SendResponseSuccess(c, nil, processStatus)
}
