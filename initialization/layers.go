package initialization

import (
	"DragDrop-Files/internal/infrastructure/repository"
	"DragDrop-Files/internal/infrastructure/repository/postgres"
	"DragDrop-Files/trash/action"
	"DragDrop-Files/trash/route"
	"DragDrop-Files/trash/service"
	"github.com/jmoiron/sqlx"
)

func InitLayers() (routes *route.Route, businessDatabase *sqlx.DB) {
	businessDatabase = postgres.NewBusinessDatabase(ConfigService)
	minioClient := repository.NewMinioStorage(ConfigService.Minio)
	sources := postgres.Sources{
		BusinessDB: businessDatabase,
	}
	repositories := postgres.NewRepository(&sources)
	domains := service.NewDomain(repositories, ConfigService, minioClient)
	actions := action.NewAction(domains)

	routes = route.NewRoute(actions)
	return routes, businessDatabase
}
