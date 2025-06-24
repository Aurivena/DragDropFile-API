package persistence

import (
	"DragDrop-Files/models"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

func NewBusinessDatabase(config *models.ConfigService) *sqlx.DB {
	fmt.Println("start database connected")
	database, err := NewPostgresDB(&PostgresDBConfig{
		Host:     config.BusinessDB.Host,
		Port:     config.BusinessDB.Port,
		Username: config.BusinessDB.Username,
		Password: config.BusinessDB.Password,
		DBName:   config.BusinessDB.DBName,
		SSLMode:  config.BusinessDB.SSLMode,
	})
	if err != nil {
		logrus.Fatalf("failed to initialize business db: %s", err.Error())
	}
	fmt.Println("database connected")
	return database
}

func NewMinioStorage(cfg models.MinioConfig) *minio.Client {
	client, err := minio.New(cfg.Endpoint, &minio.Options{Creds: credentials.NewStaticV4(cfg.User, cfg.Password, ""),
		Secure: cfg.SSL})
	if err != nil {
		return nil
	}

	return client
}
