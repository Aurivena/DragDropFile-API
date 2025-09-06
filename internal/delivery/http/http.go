package http

import (
	"DragDrop-Files/internal/application"
	"DragDrop-Files/internal/delivery/http/handler/file"
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/internal/middleware"
	"DragDrop-Files/server"
	"strings"
	"time"

	"github.com/Aurivena/spond/v2/core"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Http struct {
	File       *file.Handler
	Middleware *middleware.Middleware
}

func NewHttp(application *application.Application, spond *core.Spond, middleware *middleware.Middleware) *Http {
	return &Http{
		File:       file.New(application, spond),
		Middleware: middleware,
	}
}

func (h *Http) InitHTTPHttps(config *entity.ServerConfig) *gin.Engine {
	ginSetMode(config.ServerMode)
	gHttp := gin.Default()
	allowOrigins := strings.Split(config.Domain, ",")

	gHttp.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"X-Session-ID", "X-Password", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := gHttp.Group("/api")
	{
		file := api.Group("/file", h.Middleware.Session)
		{
			update := file.Group("update")
			{
				update.PUT("/deleted", h.File.CountDayToDeleted)
				update.PUT("/password", h.File.Password)
				update.PUT("/count-download", h.File.CountDownload)
				update.PUT("/description", h.File.Description)
			}

			file.POST("/save", h.File.Execute)
		}
		fileID := api.Group("/file/:id", h.Middleware.FileID)
		{
			fileID.GET("/data", h.File.DataFile)
			fileID.GET("", h.File.Get)
		}

	}

	return gHttp
}

func ginSetMode(serverMode string) {
	if serverMode == server.DEVELOPMENT {
		gin.SetMode(gin.ReleaseMode)
	}
}
