package file

import "github.com/Aurivena/answer"

type Update interface {
	CountDownload(count int, sessionID string) answer.ErrorCode
	DateDeleted(countDayToDeleted int, sessionID string) answer.ErrorCode
	Password(password, sessionID string) answer.ErrorCode
	Description(description, sessionID string) answer.ErrorCode
}
