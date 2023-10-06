package storage

import (
	"context"
	"keeper/internal/config"
	"keeper/internal/logger"
	"keeper/internal/model"
	"keeper/internal/utils"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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
			login:    "user13",
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
	log := logger.InitLog(logrus.InfoLevel)
	config, err := config.GetConfig(log)
	require.NoError(t, err)
	s, err := NewStorage(ctx, log, config)
	require.NoError(t, err)
	return ctx, s
}

func TestStorageInsertData(t *testing.T) {

	tests := []struct {
		name    string
		data    model.DataBlock
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Успешное добавление данных",
			data: model.DataBlock{
				DataKeyWord: "key555",
				Data:        "data",
				MetaData:    "metadata",
				Login:       "user3",
			},
			wantErr: false,
		},
		{
			name: "Добавление для несуществующего пользователя",
			data: model.DataBlock{
				DataKeyWord: "key558",
				Data:        "data",
				MetaData:    "metadata",
				Login:       "user_",
			},
			wantErr: true,
		},
	}
	ctx, s := initStorage(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.InsertData(ctx, tt.data); (err != nil) != tt.wantErr {
				t.Errorf("storage.InsertData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorageChangeData(t *testing.T) {
	tests := []struct {
		name    string
		data    model.DataBlock
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Успешное изменение",
			data: model.DataBlock{
				DataKeyWord: "key555",
				Login:       "user3",
				Data:        "changed_data",
				MetaData:    "changed_metadata",
			},
			wantErr: false,
		},
	}
	ctx, s := initStorage(t)
	secretPassword := os.Getenv("GOPRIVATE")
	require.NotEmpty(t, secretPassword)
	log := logger.InitLog(logrus.InfoLevel)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataCipher, err := utils.GCMDataCipher(tt.data.Data, secretPassword, log)
			require.NoError(t, err)
			tt.data.CipherData = dataCipher
			if err = s.ChangeData(ctx, tt.data); (err != nil) != tt.wantErr {
				t.Errorf("storage.ChangeData() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			changedData, err := s.GetData(ctx, tt.data.Login, tt.data.DataKeyWord)
			require.NoError(t, err)
			changed := changedData[0]

			dataDecipher, err := utils.GCMDataDecipher(changed.CipherData, secretPassword,
				log)
			require.NoError(t, err)
			assert.Equal(t, tt.data.Data, dataDecipher)
		})
	}
}

func TestStorageGetData(t *testing.T) {
	tests := []struct {
		name        string
		login       string
		dataKeyWord string
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name:        "Ошибка данные не выбраны",
			login:       "user_",
			dataKeyWord: "key555",
			wantErr:     true,
		},
	}
	ctx, s := initStorage(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.GetData(ctx, tt.login, tt.dataKeyWord)
			if (err != nil) != tt.wantErr {
				t.Errorf("storage.GetData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStorageDeleteData(t *testing.T) {
	tests := []struct {
		name        string
		login       string
		dataKeyWord string
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name:        "Успешное удаление",
			login:       "user3",
			dataKeyWord: "key555",
			wantErr:     false,
		},
	}
	ctx, s := initStorage(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.DeleteData(ctx, tt.login, tt.dataKeyWord); (err != nil) != tt.wantErr {
				t.Errorf("storage.DeleteData() error = %v, wantErr %v", err, tt.wantErr)
			}

			_, err := s.GetData(ctx, tt.login, tt.dataKeyWord)
			assert.Equal(t, model.ErrNoRowsSelected, err)
		})
	}
}
