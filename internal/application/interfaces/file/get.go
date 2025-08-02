package file

import (
	"DragDrop-Files/internal/domain/entity"
	"github.com/Aurivena/answer"
)

type Get interface {
	File(id, password string) (*entity.GetFileOutput, answer.ErrorCode)
	Data(id string) (*entity.DataOutput, answer.ErrorCode)
}
