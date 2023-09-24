package utils

import (
	"keeper/internal/logger"
	"keeper/internal/model"
	"os"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateJWTToken(t *testing.T) {

	tests := []struct {
		name    string
		login   string
		log     *logrus.Logger
		wantErr bool
	}{
		{
			"Успешное создание JWT токена",
			"user1",
			logger.InitLog(),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwtString, err := GenerateJWTToken(tt.login, tt.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateJWTToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.NotEmpty(t, jwtString)
			goprivate := os.Getenv("GOPRIVATE")

			tk := &model.Token{}
			token, err := jwt.ParseWithClaims(jwtString, tk, func(token *jwt.Token) (interface{}, error) {
				return []byte(goprivate), nil
			})
			require.NoError(t, err)
			require.True(t, token.Valid)
			claims, ok := token.Claims.(jwt.MapClaims)
			if ok {
				retLogin := claims["login"].(string)
				require.Equal(t, tt.login, retLogin)
			}
		})
	}
}

func TestPasswordHash(t *testing.T) {

	tests := []struct {
		name     string
		password string
	}{
		{
			"Генерация хэша",
			"123456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := PasswordHash(tt.password)
			assert.NotEmpty(t, hash)
		})
	}
}
