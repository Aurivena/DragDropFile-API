package initialization

import (
	"DragDrop-Files/models"
	"encoding/json"
	"github.com/Aurivena/answer"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	configFilePath = "/app/config.json"
)

var (
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
	logrus.Info("load local config")
	if err := loadConfig(); err != nil {
		return err
	}
	logrus.Info("load local config success")

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
