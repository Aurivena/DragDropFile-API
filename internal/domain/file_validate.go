package domain

import (
	"DragDrop-Files/internal/domain/entity"
	"time"
)

func ValidateFile(password string, file *entity.File) error {
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

func validatePassword(password string, file *entity.File) error {
	if file.Password == nil && password == "" {
		return nil
	}
	if file.Password == nil || *file.Password != password {
		return PasswordInvalidError
	}

	return nil
}

func validateDateDeleted(file *entity.File) error {
	now := time.Now().UTC()
	if !now.Before(file.TimeDeleted.UTC()) {
		return ErrFileDeleted
	}

	return nil
}

func validateCountDownload(file *entity.File) error {
	if file.CountDownload == 0 {
		return ErrFileDeleted
	}

	return nil
}
