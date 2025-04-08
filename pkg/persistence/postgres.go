package persistence

import (
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type PostgresDBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

const (
	dbDriverName = "pgx"
)

func NewPostgresDB(config *PostgresDBConfig, serverMode string) (*sqlx.DB, error) {
	db, err := getDBConnection(config)

	if err == nil {
		return db, nil
	} else {
		fmt.Println(err.Error())
	}

	return nil, err
}

func getDBConnection(config *PostgresDBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect(dbDriverName, getConnectionString(config))
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	db.SetMaxOpenConns(60)
	return db, nil
}

func getConnectionString(config *PostgresDBConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode)
}
