package file

import (
	"DragDrop-Files/internal/application"

	"github.com/Aurivena/spond/v2/core"
)

type Handler struct {
	application *application.Application
	spond       *core.Spond
}

func New(application *application.Application, spond *core.Spond) *Handler {
	return &Handler{
		application: application,
		spond:       spond,
	}
}
