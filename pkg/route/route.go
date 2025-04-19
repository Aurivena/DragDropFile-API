package route

import (
	"DragDrop-Files/model"
	"DragDrop-Files/pkg/action"
	"DragDrop-Files/server"

	"github.com/gin-gonic/gin"
)

type Route struct {
	action *action.Action
}

func NewRoute(action *action.Action) *Route {
	return &Route{action: action}
}

func (r *Route) InitHTTPRoutes(config *model.ServerConfig) *gin.Engine {
	ginSetMode(config.ServerMode)
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/save", r.SaveFile)
		api.GET("/file/:id", r.Get)
	}

	return router
}

func ginSetMode(serverMode string) {
	if serverMode == server.DEVELOPMENT {
		gin.SetMode(gin.ReleaseMode)
	}
}
