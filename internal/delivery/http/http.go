package http

import (
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/server"
	"DragDrop-Files/trash/action"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type Http struct {
	action *action.Action
}

func NewHttp(action *action.Action) *Http {
	return &Http{action: action}
}

func (h *Http) InitHTTPHttps(config *entity.ServerConfig) *gin.Engine {
	ginSetMode(config.ServerMode)
	Http := gin.Default()
	allowOrigins := strings.Split(config.Domain, ",")

	Http.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"X-Session-ID", "X-Password", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := Http.Group("/api")
	{
		api.GET("/file/:id/data", h.GetDataFile)
		api.POST("/file/save", h.SaveFile)
		api.GET("/file/:id", h.Get)
		api.PUT("/file/update/deleted", h.UpdateCountDayToDeleted)
		api.PUT("/file/update/password", h.UpdatePassword)
		api.PUT("/file/update/count-download", h.UpdateCountDownload)
		api.PUT("/file/update/description", h.UpdateDescription)
	}

	return Http
}

func ginSetMode(serverMode string) {
	if serverMode == server.DEVELOPMENT {
		gin.SetMode(gin.ReleaseMode)
	}
}
