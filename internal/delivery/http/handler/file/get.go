package file

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// @Summary      Получить данные файла
// @Description  Получает данные файла по указаному id.
// @Tags         File
// @Produce      json
// @Param        id path string true "Идентификатор файла"
// @Success      200 {object} entity.DataOutput "Файл успешно сохранен"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /file/:id/data [get]
func (h *Handler) DataFile(c *gin.Context) {
	id := c.Param("id")

	output, processStatus := h.Data(id)
	if processStatus {

	}
	h.spond.SendResponseSuccess(c.Writer, output, processStatus)
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
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	password := c.GetHeader("X-Password")

	out, err := h.application.FileGet.File(id, password)
	if err. != OK {
		h.spond.SendResponseSuccess(c, nil, processStatus)
		return
	}

	objInfo, err := out.File.Stat()
	if err != nil {
		log.Printf("Ошибка Stat() для объекта %s: %v", id, err)
		SendResponseSuccess(c, nil, InternalServerError)
		out.File.Close()
		return
	}

	contentDisposition := fmt.Sprintf("attachment; filename=\"%s\"", out.Name)

	c.DataFromReader(OK,
		objInfo.Size,
		objInfo.ContentType,
		out.File,
		map[string]string{
			"Content-Disposition": contentDisposition,
		},
	)
}
