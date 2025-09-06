package file

import (
	"time"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/sirupsen/logrus"
)

func (a *File) CountDownload(count int, sessionID string) *envelope.AppError {
	if err := a.postgresql.FileUpdate.CountDownload(count, sessionID); err != nil {
		logrus.Error(err)
		return a.InternalServerError()
	}

	return nil
}
func (a *File) DateDeleted(countDayToDeleted int, sessionID string) *envelope.AppError {
	dateDeleted := time.Now().UTC().Add(time.Hour * 24 * time.Duration(countDayToDeleted))
	if err := a.postgresql.FileUpdate.DateDeleted(dateDeleted, sessionID); err != nil {
		logrus.Error(err)
		return a.InternalServerError()
	}

	return nil
}
func (a *File) Password(password, sessionID string) *envelope.AppError {
	if err := a.postgresql.FileUpdate.Password(password, sessionID); err != nil {
		logrus.Error(err)
		return a.InternalServerError()
	}

	return nil
}

func (a *File) Description(description, sessionID string) *envelope.AppError {
	if err := a.postgresql.FileUpdate.Description(description, sessionID); err != nil {
		logrus.Error(err)
		return a.InternalServerError()
	}

	return nil
}
