package persistence

import (
	"DragDrop-Files/models"
	"fmt"

	"github.com/jmoiron/sqlx"
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
	}, config.Server.ServerMode)
	if err != nil {
		logrus.Fatalf("failed to initialize business db: %s", err.Error())
	}
	fmt.Println("database connected")
	return database
}
