package initialization

import (
	"DragDrop-Files/internal/domain/entity"
	"encoding/json"

	"os"

	"github.com/sirupsen/logrus"
)

const (
	configFilePath = "configs/config.json"
)

var (
	ConfigService = &entity.ConfigService{}
)

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

	if err = json.Unmarshal(file, &ConfigService); err != nil {
		return err
	}

	return nil
}
