package file

import (
	"fmt"
	"log"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary      Получить данные файла
// @Description  Получает данные файла по указаному id.
// @Tags         Get
// @Produce      json
// @Param        id path string true "Идентификатор файла"
// @Success      200 {object} entity.DataOutput "Файл успешно сохранен"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/:id/data [get]
func (h *Handler) DataFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		logrus.Error("missing session ID header")
		h.spond.SendResponseError(c.Writer, *h.ErrorID())
		return
	}

	output, errResp := h.application.File.Data(id)
	if errResp != nil {
		h.spond.SendResponseError(c.Writer, *errResp)
		return
	}
	h.spond.SendResponseSuccess(c.Writer, envelope.Success, output)
}

// @Summary      Получить файл
// @Description  Получает файл по переданному идентификатору.
// @Tags         Get
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
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		logrus.Error("missing session ID header")
		h.spond.SendResponseError(c.Writer, *h.ErrorID())
		return
	}

	password := c.GetHeader("X-Password")
	if password == "" {
		logrus.Error("missing session ID header")
		h.spond.SendResponseError(c.Writer, *h.ErrorPassword())
		return
	}

	out, errResp := h.application.File.Get(id, password)
	if errResp != nil {
		h.spond.SendResponseError(c.Writer, *errResp)
		return
	}

	objInfo, err := out.File.Stat()
	if err != nil {
		log.Printf("Ошибка Stat() для объекта %s: %v", id, err)
		h.spond.SendResponseError(c.Writer, *h.spond.BuildError(
			envelope.NotFound,
			"Ошибка при обработке файла",
			"Не удалось обработать файл ваш файл",
			"1. Обратитесь к администратору.",
		))
		_ = out.File.Close()
		return
	}

	contentDisposition := fmt.Sprintf("attachment; filename=\"%s\"", out.Name)

	c.DataFromReader(int(envelope.Success),
		objInfo.Size,
		objInfo.ContentType,
		out.File,
		map[string]string{
			"Content-Disposition": contentDisposition,
		},
	)
}
