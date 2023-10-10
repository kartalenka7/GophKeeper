package config

import (
	"keeper/internal/logger"
	"keeper/internal/model"
	"testing"

	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadConfigFile(t *testing.T) {
	var cfg model.Config
	err := env.Parse(&cfg)
	require.NoError(t, err)

	tests := []struct {
		name    string
		log     *logrus.Logger
		want    model.Config
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Чтение конфигурационного файла",
			log:     logger.InitLog(logrus.InfoLevel),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := readConfigFile(cfg.ConfigFile, tt.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("readConfigFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotEmpty(t, cfg)
		})
	}
}

func TestGetConfig(t *testing.T) {

	tests := []struct {
		name    string
		log     *logrus.Logger
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Получение конфигурации",
			log:     logger.InitLog(logrus.InfoLevel),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := GetConfig(tt.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotEmpty(t, cfg)
		})
	}
}
