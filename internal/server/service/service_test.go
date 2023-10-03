// Пакет service используется в качестве прослойки между
// пакетом хэндлеров и пакетом хранилища
package service

import (
	"context"
	"keeper/internal/logger"
	"keeper/internal/server/service/mocks"
	"keeper/internal/utils"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRegister(t *testing.T) {

	mockStorage := new(mocks.Storer)

	tests := []struct {
		name     string
		s        *service
		login    string
		password string
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "Успешная регистрация",
			s: &service{
				storage: mockStorage,
				log:     logger.InitLog(logrus.InfoLevel),
			},
			login:    "user8",
			password: "123456",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockStorage.On("AddUser", ctx, tt.login, utils.PasswordHash(tt.password)).Return(nil)

			jwtString, err := tt.s.UserRegister(ctx, tt.login, tt.password)

			require.NoError(t, err)
			assert.NotNil(t, jwtString)

		})
	}
}

func TestUserAuthentification(t *testing.T) {

	mockStorage := new(mocks.Storer)

	tests := []struct {
		name     string
		s        *service
		login    string
		password string
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "Успешная аутентификация",
			s: &service{
				storage: mockStorage,
				log:     logger.InitLog(logrus.InfoLevel),
			},
			login:    "user1",
			password: "123456",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockStorage.On("CheckUserAuth", ctx, tt.login, tt.password).Return(nil)

			jwtString, err := tt.s.UserAuthentification(ctx, tt.login, tt.password)

			require.NoError(t, err)
			assert.NotNil(t, jwtString)

		})
	}
}
