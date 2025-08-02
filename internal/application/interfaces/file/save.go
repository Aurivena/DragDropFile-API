package file

import (
	"DragDrop-Files/internal/domain/entity"
	"context"
	"github.com/Aurivena/answer"
	"mime/multipart"
)

type Save interface {
	Files(ctx context.Context, sessionID string, files []multipart.File, headers []*multipart.FileHeader) (*entity.FileSaveOutput, answer.ErrorCode)
}
