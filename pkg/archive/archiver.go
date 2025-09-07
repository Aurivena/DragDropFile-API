package archive

import (
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/pkg/fileops"
	"archive/zip"
	"bytes"
	"fmt"
	"io"

	_ "github.com/sirupsen/logrus"
)

func ZipFiles(files []entity.File, id string) ([]byte, error) {
	var buff bytes.Buffer
	zipW := zip.NewWriter(&buff)

	//TODO добавить worker pool
	for _, data := range files {
		fileBytes, err := fileops.DecodeFile(data.FileBase64)
		if err != nil {
			_ = zipW.Close()
			return nil, fmt.Errorf("ошибка при обработке файла %s: %w", id, err)
		}

		if len(fileBytes) == 0 {
			continue
		}
		header := &zip.FileHeader{
			Name:   data.Name,
			Method: zip.Deflate,
		}

		fileInZip, err := zipW.CreateHeader(header)
		if err != nil {
			_ = zipW.Close()
			return nil, fmt.Errorf("ошибка при создании файла %s в zip-архиве: %w", header.Name, err)
		}

		_, err = io.Copy(fileInZip, bytes.NewReader(fileBytes))
		if err != nil {
			_ = zipW.Close()
			return nil, fmt.Errorf("ошибка при записи содержимого файла %d в zip-архив: %w", header.Name, err)
		}
	}

	if err := zipW.Close(); err != nil {
		return nil, fmt.Errorf("ошибка при закрытии zip-архива: %w", err)
	}

	return buff.Bytes(), nil
}
