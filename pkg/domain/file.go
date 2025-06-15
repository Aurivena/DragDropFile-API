package domain

import (
	"DragDrop-Files/models"
	"DragDrop-Files/pkg/persistence"
	"archive/zip"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
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

func (d *FileService) Delete(id int) error {
	return d.pers.Delete(id)
}

func (d *FileService) GetByID(id string) (*models.FileOutput, error) {
	return d.pers.GetByID(id)
}

func (d *FileService) UpdateDescription(description, id string) error {
	return d.pers.UpdateDescription(description, id)
}

func (d *FileService) GetDataFile(id string) (*models.DataOutput, error) {
	out, err := d.pers.GetDataFile(id)
	if err != nil {
		if err.Error() == "no sql result" {
			return nil, errors.New("file deleted")
		}
		return nil, err
	}
	data := models.DataOutput{
		Password:      out.Password,
		DateDeleted:   out.DateDeleted,
		CountDownload: out.CountDownload,
		Description:   out.Description,
	}
	return &data, nil
}

func NewFileService(pers *persistence.Persistence) *FileService {
	return &FileService{pers: pers}
}

func (d *FileService) Create(ctx context.Context, input models.FileSave) error {
	return d.pers.Create(ctx, input)
}

func (d *FileService) GetFilesBySession(sessionID string) ([]models.FileOutput, error) {
	return d.pers.GetFilesBySessionNotZip(sessionID)
}

func (d *FileService) GetZipMetaBySession(sessionID string) (*models.FileOutput, error) {
	return d.pers.GetZipMetaBySession(sessionID)
}

func (d *FileService) DeleteFilesBySessionID(sessionID string) error {
	return d.pers.DeleteFilesBySessionID(sessionID)
}

func (d *FileService) UpdateCountDownload(count int, session string) error {
	return d.pers.File.UpdateCountDownload(count, session)
}
func (d *FileService) UpdateDateDeleted(countDayToDeleted int, session string) error {
	dateDeleted := time.Now().UTC().Add(time.Hour * 24 * time.Duration(countDayToDeleted))
	return d.pers.File.UpdateDateDeleted(dateDeleted, session)
}
func (d *FileService) UpdatePassword(password, session string) error {
	return d.pers.File.UpdatePassword(password, session)
}

func (d *FileService) ValidatePassword(input *models.FileGet) error {

	out, err := d.pers.GetByID(input.ID)
	if err != nil {
		return err
	}
	if out.Password == nil && input.Password == "" {
		return nil
	}
	if out.Password == nil || *out.Password != input.Password {
		return fmt.Errorf("пароли не совпадают")
	}

	return nil
}

func (d *FileService) ValidateDateDeleted(id string) error {
	out, err := d.pers.GetByID(id)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	if !now.Before(out.DateDeleted.UTC()) {
		if err := d.deleteFiles(id); err != nil {
			return err
		}
		return errors.New("file deleted")
	}

	return nil
}

func (d *FileService) ValidateCountDownload(id string) error {
	out, err := d.pers.GetByID(id)
	if err != nil {
		return err
	}
	if out.CountDownload == 0 {
		err := d.deleteFiles(id)
		if err != nil {
			return err
		}
		return errors.New("file deleted")
	}

	if out.CountDownload > 0 {
		c := out.CountDownload - 1
		err := d.pers.File.UpdateCountDownload(c, out.Session)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *FileService) ZipFiles(files []models.File, id string) ([]byte, error) {
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

func (d *FileService) deleteFiles(id string) error {
	out, err := d.pers.GetByID(id)
	if err != nil {
		return err
	}

	err = d.pers.DeleteFilesBySessionID(out.Session)
	if err != nil {
		return err
	}
	return nil
}
