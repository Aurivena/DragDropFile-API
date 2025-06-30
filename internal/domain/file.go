package domain

import (
	"DragDrop-Files/internal/persistence"
	"DragDrop-Files/models"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
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

func GenerateID(lenCode int) (string, error) {
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

func (d *FileService) DeleteFiles(id string) error {
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
