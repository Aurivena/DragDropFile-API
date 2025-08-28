package initialization

import (
	"DragDrop-Files/internal/application"
	"DragDrop-Files/internal/delivery/http"
	"DragDrop-Files/internal/infrastructure/minio"
	"DragDrop-Files/internal/infrastructure/repository/postgres"

	"github.com/Aurivena/spond/v2/core"
	"github.com/jmoiron/sqlx"
)

func InitLayers() (delivery *http.Http, businessDatabase *sqlx.DB) {
	spond := core.NewSpond()
	businessDatabase = postgres.NewBusinessDatabase(ConfigService)
	minioClient := minio.NewMinioStorage(ConfigService.Minio)
	sources := postgres.Sources{
		BusinessDB: businessDatabase,
	}
	repositories := postgres.New(&sources)
	minioStorage := minio.New(&ConfigService.Minio, minioClient)
	app := application.New(repositories, minioStorage, spond)
	delivery = http.NewHttp(app, spond)
	return delivery, businessDatabase
}
