package http

import (
	"DragDrop-Files/internal/application"
	"DragDrop-Files/internal/delivery/http/handler/file"
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/server"
	"strings"
	"time"

	"github.com/Aurivena/spond/v2/core"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Http struct {
	File *file.Handler
}

func NewHttp(application *application.Application, spond *core.Spond) *Http {
	return &Http{
		File: file.New(application, spond),
	}
}

func (h *Http) InitHTTPHttps(config *entity.ServerConfig) *gin.Engine {
	ginSetMode(config.ServerMode)
	gHttp := gin.Default()
	allowOrigins := strings.Split(config.Domain, ",")

	gHttp.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"X-SessionID-ID", "X-Password", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := gHttp.Group("/api")
	{
		api.GET("/file/:id/data", h.File.DataFile)
		api.POST("/file/save", h.File.Execute)
		api.GET("/file/:id", h.File.Get)
		api.PUT("/file/update/deleted", h.File.CountDayToDeleted)
		api.PUT("/file/update/password", h.File.Password)
		api.PUT("/file/update/count-download", h.File.CountDownload)
		api.PUT("/file/update/description", h.File.Description)
	}

	return gHttp
}

func ginSetMode(serverMode string) {
	if serverMode == server.DEVELOPMENT {
		gin.SetMode(gin.ReleaseMode)
	}
}
