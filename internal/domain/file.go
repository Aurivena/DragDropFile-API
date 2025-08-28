package domain

import (
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

func SetFileID(id string) (string, error) {
	if id != "" {
		return id, nil
	}

	newID, err := uuid.NewV7()
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return newID.String(), nil
}
