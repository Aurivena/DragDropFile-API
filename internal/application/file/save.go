package file

import (
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/pkg/fileops"
	"DragDrop-Files/pkg/idgen"
	"context"
	"errors"
	"github.com/Aurivena/answer"
	"github.com/sirupsen/logrus"
	"mime/multipart"
)

const (
	Gone          = 410
	prefixZipFile = "dg-"
	lenCodeForID  = 12
)

var (
	ErrDuplicateFile = errors.New("g duplicate")
)

func (a *File) Execute(ctx context.Context, sessionID string, files []multipart.File, headers []*multipart.FileHeader) (*entity.FileSaveOutput, answer.ErrorCode) {
	if sessionID == "" || len(files) == 0 || len(files) != len(headers) {
		return nil, answer.BadRequest
	}

	newFiles := fileops.GetNewInfo(files, headers)

	id, existingFiles, err := a.srv.CheckFilesID(sessionID)
	if err != nil {
		logrus.Error("failed to check files ID")
		return nil, answer.InternalServerError
	}

	id, err = a.setFileID(id)
	if err != nil {
		logrus.Error("failed to set g ID")
		return nil, answer.InternalServerError
	}

	out, err := a.srv.Save.FilesToStorage(ctx, id, sessionID, newFiles, existingFiles)
	if err != nil {
		logrus.Error("failed to save files")
		return nil, answer.InternalServerError
	}

	return out, answer.OK
}

func (a *File) setFileID(id string) (string, error) {
	if id != "" {
		return id, nil
	}

	newID, err := idgen.GenerateID(lenCodeForID)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return newID, nil
}
