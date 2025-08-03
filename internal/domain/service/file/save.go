package file

import (
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/pkg/archive"
	"DragDrop-Files/pkg/fileops"
	"DragDrop-Files/pkg/idgen"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	lenCodeForPrefix = 3
	prefixZipFile    = "dg-"
)

func (s *File) Execute(ctx context.Context, id, sessionID string, newFiles, oldFiles []entity.File) (*entity.FileSaveOutput, error) {
	var (
		wg             sync.WaitGroup
		mu             sync.Mutex
		processedFiles []entity.File
	)

	prefix, err := idgen.GenerateID(lenCodeForPrefix)
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

			if err = s.checkValidDownloadFile(ctx, data, &f, sessionID, id, prefix); err != nil {
				return
			}

			mu.Lock()
			processedFiles = append(processedFiles, f)
			mu.Unlock()
		}(file)
	}
	wg.Wait()

	processedFiles = append(processedFiles, oldFiles...)

	meta, err := s.downloadZipFile(ctx, id, sessionID, prefixZipFile, processedFiles)
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

func (s *File) downloadZipFile(ctx context.Context, id, sessionID, prefixZipFile string, files []entity.File) (*minio.UploadInfo, error) {
	fileIDZip := fmt.Sprintf("%s%s", prefixZipFile, id)
	zipData, err := archive.ZipFiles(files, fileIDZip)
	if err != nil {
		logrus.Error("failed to zip files")
		return nil, err
	}

	zipUniqueName := fmt.Sprintf("%s.zip", uuid.NewString())
	meta, err := s.downloadFile(zipData, ".zip", zipUniqueName, sessionID, fileIDZip, ctx)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

func (s *File) downloadFile(data []byte, mimeType, filename, sessionID, id string, ctx context.Context) (*minio.UploadInfo, error) {
	meta, err := s.minio.Save.File(data, sessionID, filename)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	input := entity.FileSave{
		Id:        id,
		Name:      filename,
		SessionID: sessionID,
		MimeType:  mimeType,
	}

	if err = s.repo.FileSave.Execute(ctx, input); err != nil {
		logrus.Error("failed to save g metadata")
		return nil, err
	}

	return meta, nil
}
