package file

import (
	"github.com/Aurivena/answer"
	"github.com/sirupsen/logrus"
	"time"
)

func (a *File) CountDownload(count int, sessionID string) answer.ErrorCode {
	if err := a.post.FileUpdate.CountDownload(count, sessionID); err != nil {
		logrus.Error(err)
		return answer.InternalServerError
	}

	return answer.NoContent
}
func (a *File) DateDeleted(countDayToDeleted int, sessionID string) answer.ErrorCode {
	files, err := a.post.FileGet.ZipMetaBySession(sessionID)
	if err != nil {
		logrus.Error(err)
		return answer.InternalServerError
	}
	dateDeleted := time.Now().UTC().Add(time.Hour * 24 * time.Duration(countDayToDeleted))
	if err = a.post.FileUpdate.DateDeleted(dateDeleted, files.FileID); err != nil {
		logrus.Error(err)
		return answer.InternalServerError
	}

	return answer.NoContent
}
func (a *File) Password(password, sessionID string) answer.ErrorCode {
	if err := a.post.FileUpdate.Password(password, sessionID); err != nil {
		logrus.Error(err)
		return answer.InternalServerError
	}

	return answer.NoContent
}

func (a *File) Description(description, sessionID string) answer.ErrorCode {
	if err := a.post.FileUpdate.Description(description, sessionID); err != nil {
		logrus.Error(err)
		return answer.InternalServerError
	}

	return answer.NoContent
}
