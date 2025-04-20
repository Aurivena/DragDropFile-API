package action

import (
	"DragDrop-Files/model"
	"github.com/Aurivena/answer"
)

func (a *Action) GetFile(name string) (*model.GetFileOutput, answer.ErrorCode) {
	out, err := a.domains.Minio.GetByID(name)
	if err != nil {
		return nil, answer.BadRequest
	}

	return out, answer.OK
}

func (a *Action) SaveFile(input *model.FileSave) (string, answer.ErrorCode) {
	id, err := a.domains.Minio.Save(input)
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
