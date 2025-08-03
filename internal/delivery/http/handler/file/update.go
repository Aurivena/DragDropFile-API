package file

import (
	"DragDrop-Files/internal/domain/entity"
	"github.com/Aurivena/answer"
	"github.com/gin-gonic/gin"
)

// @Summary      Обновить дату удаления файла
// @Description  Обновляет количество дней до автоматического удаления файла.
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        input body entity.DayDeletedUpdateInput true "Данные для ввода"
// @Param        X-Session-ID header string true "Идентификатор сессии пользователя"
// @Success      204 {object} string "NoContent" "Выходные данные"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/update/deleted [put]
func (h *Handler) CountDayToDeleted(c *gin.Context) {
	var input *entity.DayDeletedUpdateInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		answer.SendResponseSuccess(c, nil, answer.BadRequest)
		return
	}

	sessionID := c.GetHeader("X-Session-ID")

	processStatus := h.application.FileUpdate.DateDeleted(input.CountDayToDeleted, sessionID)
	answer.SendResponseSuccess(c, nil, processStatus)
}

// @Summary      Обновить пароль для файла
// @Description  Обновляет пароль, необходимый для доступа к файлу.
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        input body entity.PasswordUpdateInput true "Данные для ввода"
// @Param        X-Session-ID header string true "Идентификатор сессии пользователя"
// @Success      204 {object} string "NoContent" "Выходные данные"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/update/password [put]
func (h *Handler) Password(c *gin.Context) {
	var input *entity.PasswordUpdateInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		answer.SendResponseSuccess(c, nil, answer.BadRequest)
		return
	}

	sessionID := c.GetHeader("X-Session-ID")

	processStatus := h.application.FileUpdate.Password(input.Password, sessionID)
	answer.SendResponseSuccess(c, nil, processStatus)
}

// @Summary      Обновить количество загрузок файла
// @Description  Устанавливает новое ограничение по количеству скачиваний файла.
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        input body entity.CountDownloadUpdateInput true "Данные для ввода"
// @Param        X-Session-ID header string true "Идентификатор сессии пользователя"
// @Success      204 {object} string "NoContent" "Выходные данные"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/update/count-download [put]
func (h *Handler) CountDownload(c *gin.Context) {
	var input *entity.CountDownloadUpdateInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		answer.SendResponseSuccess(c, nil, answer.BadRequest)
		return
	}

	sessionID := c.GetHeader("X-Session-ID")

	processStatus := h.application.FileUpdate.CountDownload(input.CountDownload, sessionID)
	answer.SendResponseSuccess(c, nil, processStatus)
}

// @Summary      Обновить описание файла
// @Description  Устанавливает новое описание для файла.
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        input body entity.DescriptionUpdateInput true "Данные для ввода"
// @Param        X-Session-ID header string true "Идентификатор сессии пользователя"
// @Success      204 {object} string "NoContent" "Выходные данные"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/update/description [put]
func (h *Handler) Description(c *gin.Context) {
	var input *entity.DescriptionUpdateInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		answer.SendResponseSuccess(c, nil, answer.BadRequest)
		return
	}

	sessionID := c.GetHeader("X-Session-ID")

	processStatus := h.application.FileUpdate.Description(input.Description, sessionID)
	answer.SendResponseSuccess(c, nil, processStatus)
}
