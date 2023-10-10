package utils

import (
	"context"
	"keeper/internal/logger"
	"keeper/internal/model"
	"os"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
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
			goprivate := os.Getenv("GOPRIVATE")
			require.NotEmpty(t, goprivate)

			jwtString, err := GenerateJWTToken(tt.login, tt.log, goprivate)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateJWTToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.NotEmpty(t, jwtString)

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

func TestGCMDataCipher(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		log     *logrus.Logger
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Успешное шифрование",
			data:    "little gopher",
			log:     logger.InitLog(),
			wantErr: false,
		},
	}

	secretPassword := os.Getenv("GOPRIVATE")
	require.NotEmpty(t, secretPassword)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GCMDataCipher(tt.data, secretPassword, tt.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("GCMDataCipher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotEmpty(t, got)
		})
	}
}

func TestGetLoginFromContext(t *testing.T) {

	tests := []struct {
		name    string
		login   string
		ctx     context.Context
		wantErr bool
	}{
		{
			name:  "Контекст с логином в метаданных",
			login: "user1",
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.New(map[string]string{"login": "user1"})),
			wantErr: false,
		},
		{
			name:    "Контекст без метаданных",
			ctx:     context.Background(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retLogin, err := GetLoginFromContext(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLoginFromContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.login, retLogin)
		})
	}
}
