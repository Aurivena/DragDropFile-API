package file

import (
	"DragDrop-Files/internal/domain/entity"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/gin-gonic/gin"
)

// @Summary      Обновить дату удаления файла
// @Description  Обновляет количество дней до автоматического удаления файла.
// @Tags         FilePayload
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
		h.spond.SendResponseError(c.Writer, *h.ErrorParse())
		return
	}

	sessionID := c.GetHeader("X-SessionID-ID")

	if errResp := h.application.File.DateDeleted(input.CountDayToDeleted, sessionID); errResp != nil {
		h.spond.SendResponseError(c.Writer, *errResp)
		return
	}
	h.spond.SendResponseSuccess(c.Writer, envelope.NoContent, nil)
}

// @Summary      Обновить пароль для файла
// @Description  Обновляет пароль, необходимый для доступа к файлу.
// @Tags         FilePayload
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
		h.spond.SendResponseError(c.Writer, *h.ErrorParse())
		return
	}

	sessionID := c.GetHeader("X-SessionID-ID")

	if errResp := h.application.File.Password(*input.Password, sessionID); errResp != nil {
		h.spond.SendResponseError(c.Writer, *errResp)
		return
	}
	h.spond.SendResponseSuccess(c.Writer, envelope.NoContent, nil)
}

// @Summary      Обновить количество загрузок файла
// @Description  Устанавливает новое ограничение по количеству скачиваний файла.
// @Tags         FilePayload
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
		h.spond.SendResponseError(c.Writer, *h.ErrorParse())
		return
	}

	sessionID := c.GetHeader("X-SessionID-ID")

	if errResp := h.application.File.CountDownload(*input.CountDownload, sessionID); errResp != nil {
		h.spond.SendResponseError(c.Writer, *errResp)
		return
	}
	h.spond.SendResponseSuccess(c.Writer, envelope.NoContent, nil)
}

// @Summary      Обновить описание файла
// @Description  Устанавливает новое описание для файла.
// @Tags         FilePayload
// @Accept       json
// @Produce      json
// @Param        input body entity.DescriptionUpdateInput true "Данные для ввода"
// @Param        X-SessionID-ID header string true "Идентификатор сессии пользователя"
// @Success      204 {object} string "NoContent" "Выходные данные"
// @Failure      400 {object} string "Некорректные данные"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/update/description [put]
func (h *Handler) Description(c *gin.Context) {
	sessionID := c.GetHeader("X-SessionID-ID")
	if sessionID == "" {
		h.spond.SendResponseError(c.Writer, *h.ErrorSessionID())
		return
	}

	var input *entity.FileUpdateInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		h.spond.SendResponseError(c.Writer, *h.ErrorParse())
		return
	}

	if errResp := h.application.File.Description(*input.Description, sessionID); errResp != nil {
		h.spond.SendResponseError(c.Writer, *errResp)
		return
	}
	h.spond.SendResponseSuccess(c.Writer, envelope.NoContent, nil)
}
