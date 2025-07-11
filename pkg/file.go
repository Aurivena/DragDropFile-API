package pkg

import (
	"DragDrop-Files/models"
	"encoding/base64"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"strings"
	"sync"
)

func InfoNewFile(files []multipart.File, headers []*multipart.FileHeader) []models.File {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	var newFiles []models.File
	for i, file := range files {
		wg.Add(1)
		go func(f multipart.File, headers []*multipart.FileHeader, index int) {
			defer wg.Done()
			defer f.Close()

			fileData, err := getFileData(f, headers[index])
			if err != nil {
				logrus.Errorf("failed to process file %d", index)
				return
			}
			mu.Lock()

			newFiles = append(newFiles, *fileData)
			mu.Unlock()
		}(file, headers, i)
	}

	wg.Wait()

	return newFiles
}

func CheckFiles(outFile *models.GetFileOutput, file models.FileOutput, filesBase64 *[]models.File, path string) error {
	content, err := io.ReadAll(outFile.File)
	if err != nil {
		logrus.Errorf("failed to read file %s", path)

		_ = outFile.File.Close()
		return err
	}
	defer outFile.File.Close()

	encoded := base64.StdEncoding.EncodeToString(content)
	fileBase64 := fmt.Sprintf("data:%s;base64,%s", file.MimeType, encoded)

	fileInfo := models.File{
		FileBase64: fileBase64,
		Filename:   file.Name,
	}
	*filesBase64 = append(*filesBase64, fileInfo)

	return nil
}

func getFileData(file multipart.File, header *multipart.FileHeader) (*models.File, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logrus.Error("failed to read file")
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(fileBytes)
	mimeType := header.Header.Get("Content-Type")
	fileBase64 := fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)

	return &models.File{
		FileBase64: fileBase64,
		Filename:   header.Filename,
	}, nil
}
func GetMimeType(fileBase64 string) string {
	base64Data := fileBase64
	var mimeType string

	if idx := strings.Index(base64Data, ";base64,"); idx != -1 {
		parts := strings.SplitN(fileBase64, ";base64,", 2)
		if len(parts) == 2 {
			mimePart := parts[0]
			if strings.HasPrefix(mimePart, "data:") {
				mimeType = mimePart[len("data:"):]
			}
		}
	}

	return mimeType
}
func DecodeFile(fileBase64 string) ([]byte, error) {
	if !strings.HasPrefix(fileBase64, "data:") {
		logrus.Error("invalid base64 format: missing data prefix")
		return nil, fmt.Errorf("invalid base64 format: missing data prefix")
	}

	parts := strings.SplitN(fileBase64, ";base64,", 2)
	if len(parts) != 2 {
		logrus.Error("invalid base64 format: missing base64 separator")
		return nil, fmt.Errorf("invalid base64 format: missing base64 separator")
	}

	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		logrus.Error("failed to decode base64: %w", err)
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	return data, nil
}
