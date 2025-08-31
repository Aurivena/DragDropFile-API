package postgres

import (
	"DragDrop-Files/internal/domain/entity"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

type DBConfig struct {
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

func NewPostgresDB(config *DBConfig) (*sqlx.DB, error) {
	db, err := getDBConnection(config)

	if db == nil {
		logrus.Errorf("Failed to connect to postgres database: %v", err)
		return nil, err
	}

	if err = goose.SetDialect(dbDriverName); err != nil {
		logrus.Errorf("Failed to set goose dialect: %v", err)
		return nil, fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err = goose.Up(db.DB, migrationsDir); err != nil {
		return nil, err
	}

	return db, nil
}

func getDBConnection(config *DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect(dbDriverName, getConnectionString(config))
	if err != nil {
		logrus.Error(err.Error())
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(60)
	return db, nil
}

func getConnectionString(config *DBConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode)
}

func NewBusinessDatabase(config *entity.ConfigService) *sqlx.DB {
	fmt.Println("start database connected")
	database, err := NewPostgresDB(&DBConfig{
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
