package domain

import (
	"DragDrop-Files/internal/domain/entity"
	"errors"
	"fmt"
	"time"
)

func ValidateFile(password string, file *entity.FileOutput) error {
	if err := validatePassword(password, file); err != nil {
		return err
	}

	if err := validateCountDownload(file); err != nil {
		return err
	}

	if err := validateDateDeleted(file); err != nil {
		return err
	}
	return nil
}

func validatePassword(password string, file *entity.FileOutput) error {

	if file.Password == nil && password == "" {
		return nil
	}
	if file.Password == nil || *file.Password != password {
		return fmt.Errorf("пароли не совпадают")
	}

	return nil
}

func validateDateDeleted(file *entity.FileOutput) error {

	now := time.Now().UTC()
	if !now.Before(file.DateDeleted.UTC()) {
		return errors.New("file deleted")
	}

	return nil
}

func validateCountDownload(file *entity.FileOutput) error {
	if file.CountDownload == 0 {
		return errors.New("file deleted")
	}

	return nil
}
