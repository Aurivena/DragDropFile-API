package file

import (
	"github.com/Aurivena/answer"
	"github.com/sirupsen/logrus"
	"time"
)

func (a *File) UpdateCountDownload(count int, sessionID string) answer.ErrorCode {
	if err := a.repo.Update.CountDownload(count, sessionID); err != nil {
		logrus.Error(err)
		return answer.InternalServerError
	}

	return answer.NoContent
}
func (a *File) UpdateDateDeleted(countDayToDeleted int, sessionID string) answer.ErrorCode {
	files, err := a.repo.Get.ZipMetaBySession(sessionID)
	if err != nil {
		logrus.Error(err)
		return answer.InternalServerError
	}
	dateDeleted := time.Now().UTC().Add(time.Hour * 24 * time.Duration(countDayToDeleted))
	if err = a.repo.Update.DateDeleted(dateDeleted, files.FileID); err != nil {
		logrus.Error(err)
		return answer.InternalServerError
	}

	return answer.NoContent
}
func (a *File) UpdatePassword(password, sessionID string) answer.ErrorCode {
	if err := a.repo.Update.Password(password, sessionID); err != nil {
		logrus.Error(err)
		return answer.InternalServerError
	}

	return answer.NoContent
}

func (a *File) UpdateDescription(description, sessionID string) answer.ErrorCode {
	if err := a.repo.Update.Description(description, sessionID); err != nil {
		logrus.Error(err)
		return answer.InternalServerError
	}

	return answer.NoContent
}
