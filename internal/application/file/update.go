package file

import (
	"DragDrop-Files/internal/domain"
	"time"

	"github.com/sirupsen/logrus"
)

func (a *File) CountDownload(count int, sessionID string) error {
	if err := a.updater.CountDownload(count, sessionID); err != nil {
		logrus.Error(err)
		return domain.InternalError
	}

	return nil
}
func (a *File) DateDeleted(countDayToDeleted int, sessionID string) error {
	dateDeleted := time.Now().UTC().Add(time.Hour * 24 * time.Duration(countDayToDeleted))
	if err := a.updater.DateDeleted(dateDeleted, sessionID); err != nil {
		logrus.Error(err)
		return domain.InternalError
	}

	return nil
}
func (a *File) Password(password, sessionID string) error {
	if err := a.updater.Password(password, sessionID); err != nil {
		logrus.Error(err)
		return domain.InternalError
	}

	return nil
}

func (a *File) Description(description, sessionID string) error {
	if err := a.updater.Description(description, sessionID); err != nil {
		logrus.Error(err)
		return domain.InternalError
	}

	return nil
}
