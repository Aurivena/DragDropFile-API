package domain

import (
	"DragDrop-Files/models"
	"DragDrop-Files/pkg/persistence"
)

type Domain struct {
}

func NewDomain(persistence persistence.Persistence, config *models.ConfigService) *Domain {
	return &Domain{}
}
