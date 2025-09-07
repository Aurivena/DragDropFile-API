package initialization

import (
	"DragDrop-Files/internal/application"
	"DragDrop-Files/internal/delivery/http"
	"DragDrop-Files/internal/infrastructure/repository/postgres"
	"DragDrop-Files/internal/infrastructure/s3_minio"
	"DragDrop-Files/internal/middleware"

	"github.com/Aurivena/spond/v2/core"
	"github.com/jmoiron/sqlx"
)

func InitLayers() (delivery *http.Http, businessDatabase *sqlx.DB) {
	spond := core.NewSpond()
	businessDatabase = postgres.NewBusinessDatabase(ConfigService)
	minioClient := NewMinioStorage(ConfigService.Minio)
	sources := postgres.Sources{
		BusinessDB: businessDatabase,
	}
	repositories := postgres.New(&sources)
	minioStorage := s3_minio.NewS3(minioClient, &ConfigService.Minio)
	app := application.New(repositories, minioStorage)
	middleware := middleware.New(spond)
	delivery = http.NewHttp(app, spond, middleware)
	return delivery, businessDatabase
}
