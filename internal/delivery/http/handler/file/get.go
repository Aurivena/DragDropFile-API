package file

import (
	"DragDrop-Files/internal/middleware"
	"fmt"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Tags         File
// @Summary      Получить метаданные файла
// @Description  Возвращает служебные данные/метаданные для файла по его идентификатору.
// @Produce      json
// @Param        id   path   string  true  "Идентификатор файла"
// @Success      200  {object} entity.FileData  "Метаданные успешно получены"
// @Failure      404  {object} map[string]any     "Файл не найден (Spond error)"
// @Failure      500  {object} map[string]any     "Внутренняя ошибка сервера (Spond error)"
// @Router       /file/{id}/data [get]
func (h *Handler) DataFile(c *gin.Context) {
	id, _ := c.Get(middleware.CtxFileID)
	output, errResp := h.application.File.Data(id.(string))
	if errResp != nil {
		h.spond.SendResponseError(c.Writer, h.httpError(errResp))
		return
	}
	h.spond.SendResponseSuccess(c.Writer, envelope.Success, output)
}

// @Tags         File
// @Summary      Скачать файл
// @Description  Возвращает бинарное содержимое файла по идентификатору. Требует пароль в заголовке.
// @Accept       json
// @Produce      octet-stream
// @Param        id          path    string  true   "Идентификатор файла"
// @Param        X-Password  header  string  true   "Пароль для доступа к файлу"
// @Success      200         {file}  string          "Файл успешно получен"
// @Header       200         {string} Content-Disposition "attachment; filename=<имя_файла>"
// @Failure      400         {object} map[string]any "Некорректные данные (Spond error)"
// @Failure      401         {object} map[string]any "Неверный пароль (Spond error)"
// @Failure      404         {object} map[string]any "Файл не найден (Spond error)"
// @Failure      410         {object} map[string]any "Срок хранения истёк. Файл удалён (Spond error)"
// @Failure      500         {object} map[string]any "Внутренняя ошибка сервера (Spond error)"
// @Router       /file/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	password := c.GetHeader("X-Password")

	id, _ := c.Get(middleware.CtxFileID)

	out, errResp := h.application.File.Get(id.(string), password)
	if errResp != nil {
		h.spond.SendResponseError(c.Writer, h.httpError(errResp))
		return
	}

	defer out.File.Close()

	objInfo, err := out.File.Stat()
	if err != nil {
		logrus.Errorf("Ошибка Stat() для объекта %s: %v", id.(string), err)
		h.spond.SendResponseError(c.Writer, h.spond.BuildError(
			envelope.NotFound,
			"Ошибка при обработке файла",
			"Не удалось обработать файл ваш файл",
			"1. Обратитесь к администратору.",
		))
		_ = out.File.Close()
		return
	}

	contentDisposition := fmt.Sprintf("attachment; filename=%s", out.Name)

	c.DataFromReader(int(envelope.Success),
		objInfo.Size,
		objInfo.ContentType,
		out.File,
		map[string]string{
			"Content-Disposition": contentDisposition,
		},
	)
}

// @Tags         File
// @Summary      Отметить файл как зарегистрированный
// @Description  Ставит внутреннюю отметку "registered" для файла. Тело запроса не требуется.
// @Param        id   path   string  true  "Идентификатор файла"
// @Success      204  "Успешно. Тело отсутствует"
// @Failure      404  {object} map[string]any "Файл не найден (Spond error)"
// @Failure      500  {object} map[string]any "Внутренняя ошибка сервера (Spond error)"
// @Router       /file/{id}/register [post]
func (h *Handler) Registered(c *gin.Context) {
	id, _ := c.Get(middleware.CtxFileID)

	if errResp := h.application.File.Register(id.(string)); errResp != nil {
		h.spond.SendResponseError(c.Writer, h.httpError(errResp))
		return
	}

	h.spond.SendResponseSuccess(c.Writer, envelope.NoContent, nil)
}
