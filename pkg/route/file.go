package route

import (
	"DragDrop-Files/model"
	"fmt"
	"github.com/Aurivena/answer"
	"github.com/gin-gonic/gin"
	"log"
)

// @Summary      Сохранить файлы
// @Description  Сохраняет файлы и параметры, переданные пользователем.
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        input body model.FileSaveInput true "Входные данные"
// @Param        X-Session-ID header string true "Идентификатор сессии пользователя"
// @Success      200 {object} string "Файлы успешно сохранены"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/save [post]
func (r *Route) SaveFile(c *gin.Context) {
	var input *model.FileSaveInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		answer.SendResponseSuccess(c, nil, answer.BadRequest)
		return
	}

	sessionID := c.GetHeader("X-Session-ID")

	out, processStatus := r.action.Create(input, sessionID)
	if processStatus != answer.OK {
		answer.SendResponseSuccess(c, nil, answer.BadRequest)
		return
	}

	answer.SendResponseSuccess(c, out, processStatus)
}

// @Summary      Получить файл
// @Description  Получает файл по переданному идентификатору.
// @Tags         File
// @Accept       json
// @Produce      octet-stream
// @Param        id path string true "Идентификатор файла"
// @Success      200 {file} file "Файл успешно получен"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/:id [get]
func (r *Route) Get(c *gin.Context) {
	id := c.Param("id")
	var input model.FileGetInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		answer.SendResponseSuccess(c, nil, answer.BadRequest)
		return
	}

	out, processStatus := r.action.GetFile(id, &input)
	if processStatus != answer.OK {
		answer.SendResponseSuccess(c, nil, answer.BadRequest)
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
// @Param        input body model.DayDeletedUpdateInput true "Данные для ввода"
// @Param        X-Session-ID header string true "Идентификатор сессии пользователя"
// @Success      204 {object} string "NoContent" "Выходные данные"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/update/deleted [put]
func (r *Route) UpdateCountDayToDeleted(c *gin.Context) {
	var input *model.DayDeletedUpdateInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		answer.SendResponseSuccess(c, nil, answer.BadRequest)
		return
	}

	sessionID := c.GetHeader("X-Session-ID")

	processStatus := r.action.UpdateDateDeleted(input.CountDayToDeleted, sessionID)
	if processStatus != answer.NoContent {
		answer.SendResponseSuccess(c, nil, answer.BadRequest)
		return
	}

	answer.SendResponseSuccess(c, nil, processStatus)
}

// @Summary      Обновить пароль для файла
// @Description  Обновляет пароль, необходимый для доступа к файлу.
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        input body model.PasswordUpdateInput true "Данные для ввода"
// @Param        X-Session-ID header string true "Идентификатор сессии пользователя"
// @Success      204 {object} string "NoContent" "Выходные данные"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/update/password [put]
func (r *Route) UpdatePassword(c *gin.Context) {
	var input *model.PasswordUpdateInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		answer.SendResponseSuccess(c, nil, answer.BadRequest)
		return
	}

	sessionID := c.GetHeader("X-Session-ID")

	processStatus := r.action.UpdatePassword(input.Password, sessionID)
	if processStatus != answer.NoContent {
		answer.SendResponseSuccess(c, nil, answer.BadRequest)
		return
	}

	answer.SendResponseSuccess(c, nil, processStatus)
}

// @Summary      Обновить количество загрузок файла
// @Description  Устанавливает новое ограничение по количеству скачиваний файла.
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        input body model.CountDownloadUpdateInput true "Данные для ввода"
// @Param        X-Session-ID header string true "Идентификатор сессии пользователя"
// @Success      204 {object} string "NoContent" "Выходные данные"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/update/count-download [put]
func (r *Route) UpdateCountDownload(c *gin.Context) {
	var input *model.CountDownloadUpdateInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		answer.SendResponseSuccess(c, nil, answer.BadRequest)
		return
	}

	sessionID := c.GetHeader("X-Session-ID")

	processStatus := r.action.UpdateCountDownload(input.CountDownload, sessionID)
	if processStatus != answer.NoContent {
		answer.SendResponseSuccess(c, nil, answer.BadRequest)
		return
	}

	answer.SendResponseSuccess(c, nil, processStatus)
}
