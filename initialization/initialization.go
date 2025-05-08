package initialization

import (
	"DragDrop-Files/models"
	"encoding/json"
	"github.com/Aurivena/answer"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const (
	envFilePath    = `.env`
	configFilePath = "config.json"
)

var (
	Environment   = &models.Environment{}
	ConfigService = &models.ConfigService{}
)

func ErrorInitialization() error {
	err := answer.AppendCode(410, "Gone")
	if err != nil {
		return err
	}
	return nil
}
func LoadConfiguration() error {
	if err := loadEnvironment(); err != nil {
		return err
	}
	if Environment.IsReadConfig {
		logrus.Info("load local config")
		if err := loadConfig(); err != nil {
			return err
		}
		logrus.Info("load local config success")
	}

	return nil
}

func loadEnvironment() error {
	if err := godotenv.Load(envFilePath); err != nil {
		logrus.Warning("load file not found, Environment variables load from Environment")
	}
	if err := env.Parse(Environment); err != nil {
		return err
	}

	return nil
}

func loadConfig() error {
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &ConfigService)
	if err != nil {
		return err
	}

	return nil
}
