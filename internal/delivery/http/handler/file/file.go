package file

import (
	"DragDrop-Files/internal/application"
)

type Handler struct {
	application *application.Application
}

func New(application *application.Application) *Handler {
	return &Handler{application: application}
}
