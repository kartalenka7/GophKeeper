// Пакет service используется в качестве прослойки между
// пакетом хэндлеров и пакетом хранилища
package service

import (
	"context"
	"keeper/internal/logger"
	"keeper/internal/model"
	"keeper/internal/server/service/mocks"
	"keeper/internal/utils"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
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
		{
			name: "Дублирование пользователя",
			s: &service{
				storage: mockStorage,
				log:     logger.InitLog(logrus.InfoLevel),
			},
			login:    "user9",
			password: "123456",
			wantErr:  true,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.wantErr {
				mockStorage.On("AddUser", ctx, tt.login,
					utils.PasswordHash(tt.password)).Return(model.ErrUniqueViolation)
			} else {
				mockStorage.On("AddUser", ctx, tt.login,
					utils.PasswordHash(tt.password)).Return(nil)
			}
			var jwtString string
			var err error
			if jwtString, err = tt.s.UserRegister(ctx, tt.login, tt.password); (err != nil) != tt.wantErr {
				t.Errorf("storage.CheckUserAuth() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				return
			}
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
		{
			name: "Неверный пароль",
			s: &service{
				storage: mockStorage,
				log:     logger.InitLog(logrus.InfoLevel),
			},
			login:    "user1",
			password: "1234567",
			wantErr:  true,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.wantErr {
				mockStorage.On("CheckUserAuth", ctx,
					tt.login, tt.password).Return(model.ErrUserNotFound)
			} else {
				mockStorage.On("CheckUserAuth", ctx,
					tt.login, tt.password).Return(nil)
			}
			var jwtString string
			var err error

			if jwtString, err = tt.s.UserAuthentification(ctx, tt.login, tt.password); (err != nil) != tt.wantErr {
				t.Errorf("storage.CheckUserAuth() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				return
			}
			assert.NotNil(t, jwtString)

		})
	}
}

func TestServiceAddData(t *testing.T) {
	mockStorage := new(mocks.Storer)
	secretPassword := os.Getenv("GOPRIVATE")
	require.NotEmpty(t, secretPassword)

	tests := []struct {
		name          string
		s             *service
		data          model.DataBlock
		jwtStringFill bool
		wantErr       bool
	}{
		// TODO: Add test cases.
		{
			name: "Успешное добавление данных",
			s: &service{
				storage: mockStorage,
				log:     logger.InitLog(logrus.InfoLevel),
			},
			data: model.DataBlock{
				DataKeyWord: "key",
				Data:        "data",
				MetaData:    "metadata",
				Login:       "user1",
			},
			jwtStringFill: true,
			wantErr:       false,
		},
		{
			name: "В контексте нет jwt токена",
			s: &service{
				storage: mockStorage,
				log:     logger.InitLog(logrus.InfoLevel),
			},
			data: model.DataBlock{
				DataKeyWord: "key",
				Data:        "data",
				MetaData:    "metadata",
				Login:       "user1",
			},
			jwtStringFill: false,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := initContext(tt.jwtStringFill, tt.data.Login, tt.s.log,
				secretPassword)
			require.NotNil(t, ctx)

			dataCipher, err := utils.GCMDataCipher(tt.data.Data, secretPassword, tt.s.log)
			require.NoError(t, err)
			tt.data.CipherData = dataCipher

			tt.s.config.SecretPassword = secretPassword
			mockStorage.On("InsertData", ctx, tt.data).Return(nil)

			if err := tt.s.AddData(ctx, tt.data); (err != nil) != tt.wantErr {
				t.Errorf("service.AddData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceGetData(t *testing.T) {
	mockStorage := new(mocks.Storer)
	secretPassword := os.Getenv("GOPRIVATE")
	require.NotEmpty(t, secretPassword)

	tests := []struct {
		name          string
		s             *service
		login         string
		dataKeyWord   string
		jwtStringFill bool
		returnData    []model.DataBlock
		wantErr       bool
	}{
		// TODO: Add test cases.
		{
			name: "Успешное получение данных",
			s: &service{
				storage: mockStorage,
				log:     logger.InitLog(logrus.InfoLevel),
			},
			login:         "user1",
			dataKeyWord:   "key1",
			jwtStringFill: true,
			returnData: []model.DataBlock{
				{
					DataKeyWord: "key1",
					Data:        "data1",
					MetaData:    "metadata1",
				},
			},
			wantErr: false,
		},
	}

	var data []model.DataBlock
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := initContext(tt.jwtStringFill, tt.login, tt.s.log,
				secretPassword)
			require.NotNil(t, ctx)

			tt.s.config.SecretPassword = secretPassword
			cipheredData, err := utils.GCMDataCipher("data1", secretPassword, tt.s.log)
			require.NoError(t, err)

			tt.returnData[0].CipherData = cipheredData
			mockStorage.On("GetData", ctx, tt.login, tt.dataKeyWord).Return(tt.returnData, nil)

			if data, err = tt.s.GetData(ctx, tt.dataKeyWord); (err != nil) != tt.wantErr {
				t.Errorf("service.GetData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			assert.NotEmpty(t, data)
		})
	}
}

func TestServiceChangeData(t *testing.T) {
	mockStorage := new(mocks.Storer)
	secretPassword := os.Getenv("GOPRIVATE")
	require.NotEmpty(t, secretPassword)

	tests := []struct {
		name          string
		s             *service
		dataForChange model.DataBlock
		jwtStringFill bool
		wantErr       bool
	}{
		// TODO: Add test cases.
		{
			name: "Успешное изменение",
			s: &service{
				storage: mockStorage,
				log:     logger.InitLog(logrus.InfoLevel),
			},
			dataForChange: model.DataBlock{
				Login:       "user1",
				DataKeyWord: "key1",
				Data:        "data1",
				MetaData:    "metadata1",
			},
			jwtStringFill: true,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := initContext(tt.jwtStringFill, tt.dataForChange.Login, tt.s.log,
				secretPassword)
			require.NotNil(t, ctx)

			tt.s.config.SecretPassword = secretPassword
			cipheredData, err := utils.GCMDataCipher(tt.dataForChange.Data,
				secretPassword, tt.s.log)
			require.NoError(t, err)
			tt.dataForChange.CipherData = cipheredData

			mockStorage.On("ChangeData", ctx, tt.dataForChange).Return(nil)

			if err := tt.s.ChangeData(ctx, tt.dataForChange); (err != nil) != tt.wantErr {
				t.Errorf("service.ChangeData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func initContext(fillToken bool, login string, log *logrus.Logger,
	secretPassword string) context.Context {
	var ctx context.Context
	if fillToken {
		jwtString, err := utils.GenerateJWTToken(login, log, secretPassword)
		if err != nil {
			return nil
		}
		ctx = metadata.NewIncomingContext(context.Background(),
			metadata.Pairs("token", jwtString))
	} else {
		ctx = context.Background()
	}
	return ctx
}
