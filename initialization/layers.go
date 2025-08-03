package initialization

import (
	"DragDrop-Files/internal/application"
	"DragDrop-Files/internal/delivery/http"
	"DragDrop-Files/internal/domain/service"
	"DragDrop-Files/internal/infrastructure/minio"
	"DragDrop-Files/internal/infrastructure/repository/postgres"
	"github.com/jmoiron/sqlx"
)

func InitLayers() (delivery *http.Http, businessDatabase *sqlx.DB) {
	businessDatabase = postgres.NewBusinessDatabase(ConfigService)
	minioClient := minio.NewMinioStorage(ConfigService.Minio)
	sources := postgres.Sources{
		BusinessDB: businessDatabase,
	}
	repositories := postgres.New(&sources)
	minioStorage := minio.New(&ConfigService.Minio, minioClient)
	srv := service.New(repositories, minioStorage)
	application := application.New(repositories, minioStorage, srv)
	delivery = http.NewHttp(application)
	return delivery, businessDatabase
}
