package file

import (
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/pkg/fileops"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

var errFileDeleted = errors.New("file deleted")

func (s *File) File(zipFileID, id, password string) error {
	f := entity.FileGet{
		ID:       zipFileID,
		Password: password,
	}

	file, err := s.repo.FileGet.ByID(id)
	if err != nil {
		return err
	}

	if err = s.password(&f, file); err != nil {
		logrus.Error(err)
		return err
	}

	if err = s.countDownload(file); err != nil {
		logrus.Error(err)
		return s.handleIfDeleted(id, err)
	}

	if err = s.dateDeleted(file); err != nil {
		logrus.Error(err)
		return s.handleIfDeleted(id, err)
	}
	return nil
}

func (s *File) CheckFilesID(sessionID string) (string, []entity.File, error) {
	files, err := s.repo.FileGet.FilesBySessionNotZip(sessionID)
	if err != nil {
		logrus.Error("failed to g files by session")
		return "", nil, err
	}

	if len(files) == 0 {
		return "", nil, nil
	}

	var filesBase64 []entity.File
	for _, file := range files {
		path := fmt.Sprintf("%s/%s", sessionID, file.Name)
		out, err := s.minio.Get.ByFilename(path)
		if err != nil {
			logrus.Errorf("failed to g g %s from Minio", path)
			return "", nil, err
		}

		if err = fileops.CheckFiles(out, file, &filesBase64, path); err != nil {
			logrus.Error(err)
			return "", nil, err
		}

		if err = s.minio.Delete.File(file.Name); err != nil {
			logrus.Errorf("failed to delete g %s from Minio", file.Name)
			return "", nil, err
		}
	}

	if err = s.repo.FileDelete.FilesBySessionID(sessionID); err != nil {
		logrus.Error("failed to delete files by session ID")
		return "", nil, err
	}

	return files[0].FileID, filesBase64, nil
}

func (s *File) password(input *entity.FileGet, file *entity.FileOutput) error {

	if file.Password == nil && input.Password == "" {
		return nil
	}
	if file.Password == nil || *file.Password != input.Password {
		return fmt.Errorf("пароли не совпадают")
	}

	return nil
}

func (s *File) dateDeleted(file *entity.FileOutput) error {

	now := time.Now().UTC()
	if !now.Before(file.DateDeleted.UTC()) {
		return errors.New("file deleted")
	}

	return nil
}

func (s *File) countDownload(file *entity.FileOutput) error {
	if file.CountDownload == 0 {
		return errors.New("file deleted")
	}

	if file.CountDownload > 0 {
		c := file.CountDownload - 1
		err := s.repo.FileUpdate.CountDownload(c, file.Session)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *File) checkValidDownloadFile(ctx context.Context, data []byte, f *entity.File, sessionID, id, prefix string) error {
	_, err := s.downloadFile(data, fileops.GetMimeType(f.FileBase64), f.Filename, sessionID, id, ctx)
	if errors.Is(err, errFileDeleted) {
		f.Filename = fmt.Sprintf("dublicate-%s-%s", prefix, f.Filename)
		_, err = s.downloadFile(data, fileops.GetMimeType(f.FileBase64), f.Filename, sessionID, id, ctx)
		if err != nil {
			logrus.Error(err)
			return err
		}
	}

	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (s *File) handleIfDeleted(id string, err error) error {
	if errors.Is(err, errFileDeleted) {
		if err = s.minio.Delete.File(id); err != nil {
			logrus.Error(err)
			return err
		}
		return errors.New("file deleted")
	}
	return err
}
