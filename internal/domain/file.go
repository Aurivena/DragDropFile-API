package domain

import (
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/pkg/idgen"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	MimeTypeZip = ".zip"
)

func SetFileID(id string) (string, error) {
	if id != "" {
		return id, nil
	}

	newID, err := idgen.GenerateID()
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return newID, nil
}

func CheckFiles(outFile *entity.GetFileOutput, file entity.File, filesBase64 *[]entity.File, path string) error {
	content, err := io.ReadAll(outFile.File)
	if err != nil {
		logrus.Errorf("failed to read file %s", path)
		_ = outFile.File.Close()
		return err
	}
	defer outFile.File.Close()

	encoded := base64.StdEncoding.EncodeToString(content)
	fileBase64 := fmt.Sprintf("data:%s;base64,%s", file.MimeType, encoded)

	fileInfo := entity.File{
		FileBase64: fileBase64,
		Name:       file.Name,
	}
	*filesBase64 = append(*filesBase64, fileInfo)

	return nil
}

func GetNewInfo(files []multipart.File, headers []*multipart.FileHeader) []entity.File {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	var newFiles []entity.File
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

func getFileData(file multipart.File, header *multipart.FileHeader) (*entity.File, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logrus.Error("failed to read file")
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(fileBytes)
	mimeType := header.Header.Get("Content-Type")
	fileBase64 := fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)

	return &entity.File{
		FileBase64: fileBase64,
		Name:       header.Filename,
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
