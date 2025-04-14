package domain

import (
	"DragDrop-Files/model"
	"DragDrop-Files/pkg/persistence"
	"crypto/rand"
	"math/big"

	"github.com/minio/minio-go/v7"
)

type MinioService struct {
	minioClient *minio.Client
	pers        *persistence.Persistence
}

const (
	userBucket = "User"
)

func NewMinioService(minioClient *minio.Client, pers *persistence.Persistence) *MinioService {
	return &MinioService{minioClient: minioClient, pers: pers}
}

func (s *MinioService) Save(input *model.FileSave) (string, error) {
	id, err := generateID()
	if err != nil {
		return "", nil
	}

	answer, err := s.pers.File.Save(id, input.DateDeleted, input.CountDownload, input.CountDiscoveries, input.CountDay)
	if err != nil {
		return "", err
	}

	if !answer {
		return "", err
	}

	return id, nil
}

func (s *MinioService) Delete(id string) error {
	return s.pers.Delete(id)
}

func (s *MinioService) Get(id string) (*model.File, error) {
	return s.pers.Get(id)
}

func generateID() (string, error) {
	lenCode := 12
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, lenCode)
	max := big.NewInt(int64(len(letters)))
	for i := range code {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		code[i] = letters[num.Int64()]
	}

	return string(code), nil
}
