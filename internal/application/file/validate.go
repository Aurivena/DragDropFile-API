package file

import (
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/pkg/fileops"
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

func (a *File) checkFilesID(sessionID string) (string, []entity.FileFFF, error) {
	files, err := a.postgresql.FileGet.FilesBySessionNotZip(sessionID)
	if err != nil {
		logrus.Error("failed to g files by session")
		return "", nil, err
	}

	if len(files) == 0 {
		return "", nil, nil
	}

	var filesBase64 []entity.FileFFF
	for _, file := range files {
		path := fmt.Sprintf("%s/%s", sessionID, file.Name)
		out, err := a.minioStorage.Get.ByFilename(path)
		if err != nil {
			logrus.Errorf("failed to g g %s from Minio", path)
			return "", nil, err
		}

		if err = fileops.CheckFiles(out, file, &filesBase64, path); err != nil {
			logrus.Error(err)
			return "", nil, err
		}

		if err = a.minioStorage.Delete.File(file.Name); err != nil {
			logrus.Errorf("failed to delete g %s from Minio", file.Name)
			return "", nil, err
		}
	}

	if err = a.postgresql.FileDelete.FilesBySessionID(sessionID); err != nil {
		logrus.Error("failed to delete files by session ID")
		return "", nil, err
	}

	return files[0].FileID, filesBase64, nil
}

func (a *File) registerDownload(countDownload int, session string) error {
	if countDownload == 0 {
		return errors.New("file deleted")
	}

	if countDownload > 0 {
		c := countDownload - 1
		err := a.postgresql.FileUpdate.CountDownload(c, session)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *File) validDownloadFile(ctx context.Context, data []byte, f *entity.FileFFF, sessionID, id, prefix string) bool {
	_, err := a.downloadFile(ctx, data, fileops.GetMimeType(f.FileBase64), f.Filename, sessionID, id)
	if errors.Is(err, errFileDeleted) {
		f.Filename = fmt.Sprintf("dublicate-%s-%s", prefix, f.Filename)
		_, err = a.downloadFile(ctx, data, fileops.GetMimeType(f.FileBase64), f.Filename, sessionID, id)
		if err != nil {
			logrus.Error(err)
			return false
		}
	}
	if err != nil {
		logrus.Error(err)
		return false
	}
	return true
}
