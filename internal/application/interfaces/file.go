package interfaces

import (
	"DragDrop-Files/internal/domain/entity"
	"context"
	"github.com/Aurivena/answer"
	"mime/multipart"
)

type FileSave interface {
	Execute(ctx context.Context, sessionID string, files []multipart.File, headers []*multipart.FileHeader) (*entity.FileSaveOutput, answer.ErrorCode)
}

type FileGet interface {
	File(id, password string) (*entity.GetFileOutput, answer.ErrorCode)
	Data(id string) (*entity.DataOutput, answer.ErrorCode)
}

type FileUpdate interface {
	CountDownload(count int, sessionID string) answer.ErrorCode
	DateDeleted(countDayToDeleted int, sessionID string) answer.ErrorCode
	Password(password, sessionID string) answer.ErrorCode
	Description(description, sessionID string) answer.ErrorCode
}
