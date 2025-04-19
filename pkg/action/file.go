package action

import (
	"DragDrop-Files/model"
	"github.com/minio/minio-go/v7"

	"github.com/Aurivena/answer"
)

func (a *Action) GetFile(id string) (*minio.Object, answer.ErrorCode) {
	out, err := a.domains.Minio.GetByID(id)
	if err != nil {
		return nil, answer.BadRequest
	}

	return out, answer.OK
}

func (a *Action) SaveFile(file *model.FileSave) (string, answer.ErrorCode) {
	id, err := a.domains.Minio.Save(file)
	if err != nil {
		return "", answer.BadRequest
	}

	return id, answer.OK
}

func (a *Action) DeleteFile(id string) answer.ErrorCode {
	if err := a.domains.Minio.Delete(id); err != nil {
		return answer.BadRequest
	}
	return answer.NoContent
}
