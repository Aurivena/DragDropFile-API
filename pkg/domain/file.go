package domain

import (
	"DragDrop-Files/model"
	"DragDrop-Files/pkg/persistence"
	"archive/zip"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"math/big"
	"strings"
	"time"
)

type FileService struct {
	pers *persistence.Persistence
}

func NewFileService(pers *persistence.Persistence) *FileService {
	return &FileService{pers: pers}
}

func (d *FileService) Create(input model.FileSave) error {
	return d.pers.Create(input)
}

func (d *FileService) GetIdFileBySession(sessionID string) ([]string, error) {
	return d.pers.GetIdFileBySession(sessionID)
}

func (d *FileService) GetFileBySession(sessionID string) ([]model.FileOutput, error) {
	return d.pers.GetFileBySession(sessionID)
}

func (d *FileService) GetNameByID(id string) (string, error) {
	return d.pers.GetNameByID(id)
}

func (d *FileService) GetZipMetaBySession(sessionID string) (*model.FileOutput, error) {
	return d.pers.GetZipMetaBySession(sessionID)
}

func (d *FileService) Delete(id string) error {
	return d.pers.Delete(id)
}

func (d *FileService) GetMimeTypeByID(id string) (string, error) {
	return d.pers.GetMimeTypeByID(id)
}

func (d *FileService) DeleteFilesBySessionID(sessionID string) error {
	return d.pers.DeleteFilesBySessionID(sessionID)
}

func (d *FileService) GetSessionByID(id string) (string, error) {
	return d.pers.GetSessionByID(id)
}

func (d *FileService) UpdateCountDownload(count int, sessionID string) error {
	return d.pers.File.UpdateCountDownload(count, sessionID)
}
func (d *FileService) UpdateDateDeleted(countDayToDeleted int, sessionID string) error {
	dateDeleted := time.Now().UTC().Add(time.Hour * 24 * time.Duration(countDayToDeleted))
	return d.pers.File.UpdateDateDeleted(dateDeleted, sessionID)
}
func (d *FileService) UpdatePassword(password, sessionID string) error {
	return d.pers.File.UpdatePassword(password, sessionID)
}

func (d *FileService) ValidatePassword(input *model.FileGet) error {
	out, err := d.pers.Get(input.SessionID)
	if err != nil {
		return err
	}

	if out.Password != nil && *out.Password != *input.Password {
		return fmt.Errorf("пароли не совпадают")
	}

	return nil
}

func (d *FileService) ValidateDateDeleted(sessionID string) error {
	out, err := d.pers.Get(sessionID)
	if err != nil {
		return err
	}

	if out.DateDeleted != nil {
		now := time.Now().UTC()
		if !now.Before(out.DateDeleted.UTC()) {
			if err := d.pers.File.DeleteFilesBySessionID(sessionID); err != nil {
				return err
			}
			return fmt.Errorf("срок хранения этого файла истек. свяжитесь с владельцем")
		}
	}

	return nil
}
func (d *FileService) ValidateCountDownload(sessionID string) error {
	out, err := d.pers.Get(sessionID)
	if err != nil {
		return err
	}

	if out.CountDownload != nil && *out.CountDownload <= 0 {
		err := d.pers.File.DeleteFilesBySessionID(sessionID)
		if err != nil {
			return err
		}
		return fmt.Errorf("количество загрузрк для этого файла исерпано. свяжитесь с владельцем")
	}

	if out.CountDownload != nil && *out.CountDownload >= 0 {
		c := *out.CountDownload - 1
		err := d.pers.File.UpdateCountDownload(c, sessionID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *FileService) ZipFiles(files []model.File, id string) ([]byte, error) {
	var buff bytes.Buffer
	zipW := zip.NewWriter(&buff)

	for i, data := range files {
		fileBytes, err := DecodeFile(data.FileBase64)
		if err != nil {
			_ = zipW.Close()
			return nil, fmt.Errorf("ошибка при обработке файла %s: %w", id, err)
		}

		if len(fileBytes) == 0 {
			log.Printf("[zipFiles] Пустой файл %d. Пропускаем.", i)
			continue
		}
		headerName := fmt.Sprintf("%d-%s", i+1, data.Filename)
		header := &zip.FileHeader{
			Name:   headerName,
			Method: zip.Deflate,
		}

		fileInZip, err := zipW.CreateHeader(header)
		if err != nil {
			_ = zipW.Close()
			return nil, fmt.Errorf("ошибка при создании файла %d в zip-архиве: %w", header.Name, err)
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
func GenerateID() (string, error) {
	lenCode := 12
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, lenCode)
	newInt := big.NewInt(int64(len(letters)))
	for i := range code {
		num, err := rand.Int(rand.Reader, newInt)
		if err != nil {
			log.Printf("не удалось сгенерировать часть ID: %w", err)
			return "", fmt.Errorf("не удалось сгенерировать часть ID: %w", err)
		}
		code[i] = letters[num.Int64()]
	}

	return string(code), nil
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
	base64Data := fileBase64

	if idx := strings.Index(base64Data, ";base64,"); idx != -1 {
		parts := strings.SplitN(fileBase64, ";base64,", 2)
		base64Data = parts[1]
	}

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		log.Printf("Ошибка декодирования Base64 для строки '%s': %v", fileBase64[:min(len(fileBase64), 50)], err)
		return nil, fmt.Errorf("некорректные Base64 данные: %w", err)
	}

	return data, nil
}
