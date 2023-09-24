package storage

import (
	"context"
	"keeper/internal/config"
	"keeper/internal/logger"
	"keeper/internal/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorageAddUser(t *testing.T) {
	tests := []struct {
		name     string
		login    string
		password string
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name:     "Успешное добавление пользователя",
			login:    "user5",
			password: "123456",
			wantErr:  false,
		},
		{
			name:     "Добавление уже существующего пользователя",
			login:    "user4",
			password: "123456",
			wantErr:  true,
		},
	}
	ctx, s := initStorage(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.AddUser(ctx, tt.login, utils.PasswordHash(tt.password)); (err != nil) != tt.wantErr {
				t.Errorf("storage.AddUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorageCheckUserAuth(t *testing.T) {
	tests := []struct {
		name     string
		login    string
		password string
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name:     "Успешная проверка",
			login:    "user4",
			password: "123456",
			wantErr:  false,
		},
		{
			name:     "Неверный пароль",
			login:    "user4",
			password: "98765",
			wantErr:  true,
		},
	}
	ctx, s := initStorage(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.CheckUserAuth(ctx, tt.login, tt.password); (err != nil) != tt.wantErr {
				t.Errorf("storage.CheckUserAuth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func initStorage(t *testing.T) (context.Context, *storage) {
	ctx := context.Background()
	log := logger.InitLog()
	config, err := config.GetConfig(log)
	require.NoError(t, err)
	s, err := NewStorage(ctx, log, config)
	require.NoError(t, err)
	return ctx, s
}
