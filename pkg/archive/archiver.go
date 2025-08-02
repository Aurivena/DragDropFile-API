package archive

import (
	"DragDrop-Files/models"
	"DragDrop-Files/pkg/file"
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
)

func ZipFiles(files []models.File, id string) ([]byte, error) {
	var buff bytes.Buffer
	zipW := zip.NewWriter(&buff)

	for i, data := range files {
		fileBytes, err := fileops.DecodeFile(data.FileBase64)
		if err != nil {
			_ = zipW.Close()
			return nil, fmt.Errorf("ошибка при обработке файла %s: %w", id, err)
		}

		if len(fileBytes) == 0 {
			log.Printf("[zipFiles] Пустой файл %d. Пропускаем.", i)
			continue
		}
		header := &zip.FileHeader{
			Name:   data.Filename,
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

	err := zipW.Close()
	if err != nil {
		return nil, fmt.Errorf("ошибка при закрытии zip-архива: %w", err)
	}

	return buff.Bytes(), nil
}
