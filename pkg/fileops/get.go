package fileops

import (
	"DragDrop-Files/internal/domain/entity"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

func GetNewInfo(files []multipart.File, headers []*multipart.FileHeader) []entity.FileFFF {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	var newFiles []entity.FileFFF
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

func getFileData(file multipart.File, header *multipart.FileHeader) (*entity.FileFFF, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logrus.Error("failed to read file")
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(fileBytes)
	mimeType := header.Header.Get("Content-Type")
	fileBase64 := fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)

	return &entity.FileFFF{
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
