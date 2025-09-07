package file

import (
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/internal/middleware"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/gin-gonic/gin"
)

// @Tags         File
// @Summary      Обновить дату удаления файла
// @Description  Обновляет количество дней до автоматического удаления файла.
// @Accept       json
// @Produce      json
// @Param        X-Session-ID header string true "Идентификатор сессии пользователя"
// @Param        input body entity.FileUpdateInput true "Данные для ввода countDayToDeleted"
// @Success      204 "Успешно. Тело отсутствует"
// @Failure      400 {object} map[string]any "Некорректные данные (Spond error)"
// @Failure      500 {object} map[string]any "Внутренняя ошибка сервера (Spond error)"
// @Router       /file/update/deleted [put]
func (h *Handler) CountDayToDeleted(c *gin.Context) {
	var input *entity.FileUpdateInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		h.spond.SendResponseError(c.Writer, h.ErrorParse())
		return
	}

	if errResp := h.application.File.DateDeleted(input.CountDayToDeleted, c.GetHeader(middleware.Session)); errResp != nil {
		h.spond.SendResponseError(c.Writer, h.httpError(errResp))
		return
	}
	h.spond.SendResponseSuccess(c.Writer, envelope.NoContent, nil)
}

// @Tags         File
// @Summary      Обновить пароль для файла
// @Description  Обновляет пароль, необходимый для доступа к файлу.
// @Accept       json
// @Produce      json
// @Param        X-Session-ID header string true "Идентификатор сессии пользователя"
// @Param        input body entity.FileUpdateInput true "Данные для ввода password"
// @Success      204 "Успешно. Тело отсутствует"
// @Failure      400 {object} map[string]any "Некорректные данные (Spond error)"
// @Failure      500 {object} map[string]any "Внутренняя ошибка сервера (Spond error)"
// @Router       /file/update/password [put]
func (h *Handler) Password(c *gin.Context) {
	var input *entity.FileUpdateInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		h.spond.SendResponseError(c.Writer, h.ErrorParse())
		return
	}

	if errResp := h.application.File.Password(*input.Password, c.GetHeader(middleware.Session)); errResp != nil {
		h.spond.SendResponseError(c.Writer, h.httpError(errResp))
		return
	}
	h.spond.SendResponseSuccess(c.Writer, envelope.NoContent, nil)
}

// @Tags         File
// @Summary      Обновить количество загрузок файла
// @Description  Устанавливает новое ограничение по количеству скачиваний файла.
// @Accept       json
// @Produce      json
// @Param        X-Session-ID header string true "Идентификатор сессии пользователя"
// @Param        input body entity.FileUpdateInput true "Данные для ввода countDownload"
// @Success      204 "Успешно. Тело отсутствует"
// @Failure      400 {object} map[string]any "Некорректные данные (Spond error)"
// @Failure      500 {object} map[string]any "Внутренняя ошибка сервера (Spond error)"
// @Router       /file/update/count-download [put]
func (h *Handler) CountDownload(c *gin.Context) {
	var input *entity.FileUpdateInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		h.spond.SendResponseError(c.Writer, h.ErrorParse())
		return
	}

	if errResp := h.application.File.CountDownload(*input.CountDownload, c.GetHeader(middleware.Session)); errResp != nil {
		h.spond.SendResponseError(c.Writer, h.httpError(errResp))
		return
	}
	h.spond.SendResponseSuccess(c.Writer, envelope.NoContent, nil)
}

// @Tags         File
// @Summary      Обновить описание файла
// @Description  Устанавливает новое описание для файла.
// @Accept       json
// @Produce      json
// @Param        X-Session-ID header string true "Идентификатор сессии пользователя"
// @Param        input body entity.FileUpdateInput true "Данные для ввода description"
// @Success      204 "Успешно. Тело отсутствует"
// @Failure      400 {object} map[string]any "Некорректные данные (Spond error)"
// @Failure      500 {object} map[string]any "Внутренняя ошибка сервера (Spond error)"
// @Router       /file/update/description [put]
func (h *Handler) Description(c *gin.Context) {
	var input *entity.FileUpdateInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		h.spond.SendResponseError(c.Writer, h.ErrorParse())
		return
	}

	if errResp := h.application.File.Description(*input.Description, c.GetHeader(middleware.Session)); errResp != nil {
		h.spond.SendResponseError(c.Writer, h.httpError(errResp))
		return
	}
	h.spond.SendResponseSuccess(c.Writer, envelope.NoContent, nil)
}
