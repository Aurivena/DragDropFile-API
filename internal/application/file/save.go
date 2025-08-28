package file

import (
	"DragDrop-Files/internal/domain"
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/pkg/fileops"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"sync"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	prefixZipFile = "dg-"
)

var (
	ErrDuplicateFile = errors.New("file duplicate")
)

func (a *File) Execute(ctx context.Context, sessionID string, files []multipart.File, headers []*multipart.FileHeader) (*entity.FileSaveOutput, *envelope.AppError) {
	if sessionID == "" || len(files) == 0 || len(files) != len(headers) {
		return nil, a.BadRequest("1. Ваша сессия недействительна\n" + "2. Длина загруженных файлов == 0")
	}

	newFiles := fileops.GetNewInfo(files, headers)

	id, existingFiles, err := a.checkFilesID(sessionID)
	if err != nil {
		logrus.Error("failed to check files ID")
		return nil, a.InternalServerError()
	}

	id, err = domain.SetFileID(id)
	if err != nil {
		logrus.Error("failed to set g ID")
		return nil, a.InternalServerError()
	}

	out, err := a.execute(ctx, id, sessionID, newFiles, existingFiles)
	if err != nil {
		logrus.Error("failed to save files")
		return nil, a.InternalServerError()
	}

	return out, nil
}

func (a *File) execute(ctx context.Context, id, sessionID string, newFiles, oldFiles []entity.File) (*entity.FileSaveOutput, error) {
	var (
		wg             sync.WaitGroup
		mu             sync.Mutex
		processedFiles []entity.File
	)

	prefix, err := uuid.NewV7()
	if err != nil {
		logrus.Errorf("failed to generate prefix: %v", err)
		return nil, fmt.Errorf("failed to generate prefix: %w", err)
	}

	for _, file := range newFiles {
		wg.Add(1)
		go func(f entity.File) {
			defer wg.Done()

			data, err := fileops.DecodeFile(f.FileBase64)
			if err != nil {
				return
			}

			if !a.validDownloadFile(ctx, data, &f, sessionID, id, prefix.String()) {
				return
			}

			mu.Lock()
			processedFiles = append(processedFiles, f)
			mu.Unlock()
		}(file)
	}
	wg.Wait()

	processedFiles = append(processedFiles, oldFiles...)

	meta, err := a.downloadZipFile(ctx, id, sessionID, prefixZipFile, processedFiles)
	if err != nil {
		logrus.Errorf("failed to create zip g: %v", err)
		return nil, fmt.Errorf("failed to create zip g: %w", err)
	}

	return &entity.FileSaveOutput{
		ID:    id,
		Size:  meta.Size,
		Count: len(processedFiles),
	}, nil
}
