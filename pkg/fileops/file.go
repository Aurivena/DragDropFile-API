package fileops

import (
	"DragDrop-Files/internal/domain/entity"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

func CheckFiles(outFile *entity.GetFileOutput, file entity.FileOutput, filesBase64 *[]entity.File, path string) error {
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
		Filename:   file.Name,
	}
	*filesBase64 = append(*filesBase64, fileInfo)

	return nil
}
