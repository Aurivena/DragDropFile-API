package application

import (
	"DragDrop-Files/internal/domain/interfaces/repository"
	"DragDrop-Files/internal/domain/interfaces/s3"
	"DragDrop-Files/internal/domain/interfaces/service"
)

type Application struct {
}

func New(repo repository.Repository, s3 s3.S3, service service.Service) *Application {

}
