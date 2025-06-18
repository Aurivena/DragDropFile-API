package route

import (
	"DragDrop-Files/models"
	"DragDrop-Files/pkg/action"
	"DragDrop-Files/server"
	"github.com/gin-contrib/cors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Route struct {
	action *action.Action
}

func NewRoute(action *action.Action) *Route {
	return &Route{action: action}
}

func (r *Route) InitHTTPRoutes(config *models.ServerConfig) *gin.Engine {
	ginSetMode(config.ServerMode)
	router := gin.Default()
	allowOrigins := strings.Split(config.Domain, ",")

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"X-Session-ID", "X-Password", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api")
	{
		api.GET("/file/:id/data", r.GetDataFile)
		api.POST("/file/save", r.SaveFile)
		api.GET("/file/:id", r.Get)
		api.PUT("/file/update/deleted", r.UpdateCountDayToDeleted)
		api.PUT("/file/update/password", r.UpdatePassword)
		api.PUT("/file/update/count-download", r.UpdateCountDownload)
		api.PUT("/file/update/description", r.UpdateDescription)
	}

	return router
}

func ginSetMode(serverMode string) {
	if serverMode == server.DEVELOPMENT {
		gin.SetMode(gin.ReleaseMode)
	}
}
