package middleware

import (
	"github.com/Aurivena/spond/v2/envelope"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	Session = "X-Session-ID"
)

func (m *Middleware) Session(c *gin.Context) {
	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		logrus.Error("missing session ID header")
		m.spond.BuildError(
			envelope.BadRequest,
			"Отсутствует sessionID",
			"Не удалось определить sessionID пользователя",
			"1. Перезагрузите страницу.\n"+
				"2. Обратитесь к создателю ресурса.",
		)
		return
	}
}

const CtxFileID = "fileID"

func (m *Middleware) FileID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		logrus.Error("missing session ID header")
		m.spond.BuildError(
			envelope.BadRequest,
			"Проблема с ID файла",
			"Не удалось определить ID файла",
			"1. Сделайте повторно запрос.",
		)
		return
	}

	c.Set(CtxFileID, id)

}
