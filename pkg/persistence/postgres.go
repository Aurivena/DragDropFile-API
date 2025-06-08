package persistence

import (
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
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
	dbDriverName  = "pgx"
	migrationsDir = "migrations"
)

func NewPostgresDB(config *PostgresDBConfig) (*sqlx.DB, error) {
	db, err := getDBConnection(config)

	if err = goose.SetDialect(dbDriverName); err != nil {
		return nil, err
	}

	if err = goose.Up(db.DB, migrationsDir); err != nil {
		return nil, err
	}

	return db, nil
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
