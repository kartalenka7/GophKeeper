package config

import (
	"encoding/json"
	"keeper/internal/model"
	"os"

	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

// GetConfig возвращает конфигурацию приложения
func GetConfig(log *logrus.Logger) (model.Config, error) {
	var cfg model.Config

	if err := env.Parse(&cfg); err != nil {
		log.Error(err.Error())
		return cfg, err
	}
	config, err := readConfigFile(cfg.ConfigFile, log)
	if err != nil {
		log.Error(err.Error())
	}
	return config, err
}

// readConfigFile читает конфигурационный файл в формате json
func readConfigFile(filename string, log *logrus.Logger) (model.Config, error) {
	var config model.Config

	file, err := os.OpenFile(filename, os.O_RDONLY, 0664)
	if err != nil {
		log.Error(err.Error())
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		log.Error(err.Error())
	}
	return config, err
}
