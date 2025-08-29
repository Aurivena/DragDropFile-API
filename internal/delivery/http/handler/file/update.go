package file

import (
	"DragDrop-Files/internal/domain/entity"
	"github.com/gin-gonic/gin"
)

// @Summary      Обновить дату удаления файла
// @Description  Обновляет количество дней до автоматического удаления файла.
// @Tags         FileFFF
// @Accept       json
// @Produce      json
// @Param        input body entity.DayDeletedUpdateInput true "Данные для ввода"
// @Param        X-SessionID-ID header string true "Идентификатор сессии пользователя"
// @Success      204 {object} string "NoContent" "Выходные данные"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/update/deleted [put]
func (h *Handler) CountDayToDeleted(c *gin.Context) {
	var input *entity.FileUpdateInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		SendResponseSuccess(c, nil, BadRequest)
		return
	}

	sessionID := c.GetHeader("X-SessionID-ID")

	processStatus := h.application.FileUpdate.DateDeleted(input.CountDayToDeleted, sessionID)
	SendResponseSuccess(c, nil, processStatus)
}

// @Summary      Обновить пароль для файла
// @Description  Обновляет пароль, необходимый для доступа к файлу.
// @Tags         FileFFF
// @Accept       json
// @Produce      json
// @Param        input body entity.PasswordUpdateInput true "Данные для ввода"
// @Param        X-SessionID-ID header string true "Идентификатор сессии пользователя"
// @Success      204 {object} string "NoContent" "Выходные данные"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/update/password [put]
func (h *Handler) Password(c *gin.Context) {
	var input *entity.FileUpdateInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		SendResponseSuccess(c, nil, BadRequest)
		return
	}

	sessionID := c.GetHeader("X-SessionID-ID")

	processStatus := h.FileUpdate.Password(input.Password, sessionID)
	SendResponseSuccess(c, nil, processStatus)
}

// @Summary      Обновить количество загрузок файла
// @Description  Устанавливает новое ограничение по количеству скачиваний файла.
// @Tags         FileFFF
// @Accept       json
// @Produce      json
// @Param        input body entity.CountDownloadUpdateInput true "Данные для ввода"
// @Param        X-SessionID-ID header string true "Идентификатор сессии пользователя"
// @Success      204 {object} string "NoContent" "Выходные данные"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/update/count-download [put]
func (h *Handler) CountDownload(c *gin.Context) {
	var input *entity.FileUpdateInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		SendResponseSuccess(c, nil, BadRequest)
		return
	}

	sessionID := c.GetHeader("X-SessionID-ID")

	processStatus := h.application.FileUpdate.CountDownload(input.CountDownload, sessionID)
	SendResponseSuccess(c, nil, processStatus)
}

// @Summary      Обновить описание файла
// @Description  Устанавливает новое описание для файла.
// @Tags         FileFFF
// @Accept       json
// @Produce      json
// @Param        input body entity.DescriptionUpdateInput true "Данные для ввода"
// @Param        X-SessionID-ID header string true "Идентификатор сессии пользователя"
// @Success      204 {object} string "NoContent" "Выходные данные"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/update/description [put]
func (h *Handler) Description(c *gin.Context) {
	var input *entity.FileUpdateInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		SendResponseSuccess(c, nil, BadRequest)
		return
	}

	sessionID := c.GetHeader("X-SessionID-ID")

	processStatus := h.application.FileUpdate.Description(input.Description, sessionID)
	SendResponseSuccess(c, nil, processStatus)
}
