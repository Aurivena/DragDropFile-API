package service

import (
	"DragDrop-Files/internal/domain/interfaces/service"
	"DragDrop-Files/internal/infrastructure/minio"
	"DragDrop-Files/internal/infrastructure/repository/postgres"
)

type Service struct {
	service.Validate
	service.Save
}

func New(repo *postgres.Repository, minio *minio.Minio) *Service {
	return &Service{
		Validate: New(repo, minio),
		Save:     New(repo, minio),
	}
}
