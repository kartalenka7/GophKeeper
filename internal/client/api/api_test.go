package api

import (
	"context"
	"fmt"
	"io"
	"keeper/internal/client/api/mocks"
	"keeper/internal/logger"
	"keeper/internal/model"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestApiRegister(t *testing.T) {
	type args struct {
		ctx      context.Context
		log      *logrus.Logger
		service  *mocks.Service
		login    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Регистрация",
			args: args{
				ctx:      context.Background(),
				log:      logger.InitLog(logrus.InfoLevel),
				service:  new(mocks.Service),
				login:    "testlogin",
				password: "testpassword",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Создаем оригинальный os.Stdin и сохраняем его
			originalStdin := os.Stdin

			// Создаем pipe для ввода
			r, w, _ := os.Pipe()

			// Перенаправляем os.Stdin на pipe
			os.Stdin = r
			defer func() {
				os.Stdin = originalStdin
			}()

			// Создаем буфер для вывода данных (для обработки fmt.Println)
			var outputBuf strings.Builder
			wMock := io.MultiWriter(w, io.Writer(&outputBuf))

			// Вызываем функцию register с мок-объектом
			go func() {
				_, err := fmt.Fprintln(wMock, tt.args.login)
				assert.NoError(t, err)
				_, err = fmt.Fprintln(wMock, tt.args.password)
				assert.NoError(t, err)
			}()

			tt.args.service.On("Register", tt.args.ctx,
				tt.args.login, tt.args.password).Return("token", nil)

			got, err := register(tt.args.ctx, tt.args.log, tt.args.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.NotEmpty(t, got)
		})
	}
}

func TestApiAuth(t *testing.T) {
	type args struct {
		ctx      context.Context
		log      *logrus.Logger
		service  *mocks.Service
		login    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Аутентификация",
			args: args{
				ctx:      context.Background(),
				log:      logger.InitLog(logrus.InfoLevel),
				service:  new(mocks.Service),
				login:    "testlogin",
				password: "testpassword",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Создаем оригинальный os.Stdin и сохраняем его
			originalStdin := os.Stdin

			// Создаем pipe для ввода
			r, w, _ := os.Pipe()

			// Перенаправляем os.Stdin на pipe
			os.Stdin = r
			defer func() {
				os.Stdin = originalStdin
			}()

			// Создаем буфер для вывода данных (для обработки fmt.Println)
			var outputBuf strings.Builder
			wMock := io.MultiWriter(w, io.Writer(&outputBuf))

			go func() {
				_, err := fmt.Fprintln(wMock, tt.args.login)
				assert.NoError(t, err)
				_, err = fmt.Fprintln(wMock, tt.args.password)
				assert.NoError(t, err)
			}()

			tt.args.service.On("Auth", tt.args.ctx,
				tt.args.login, tt.args.password).Return("token", nil)

			got, err := auth(tt.args.ctx, tt.args.log, tt.args.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("auth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotEmpty(t, got)
		})
	}
}

func TestApiAdd(t *testing.T) {
	type args struct {
		ctx         context.Context
		log         *logrus.Logger
		service     *mocks.Service
		jwtToken    string
		dataKeyWord string
		data        string
		metadata    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Добавление данных",
			args: args{
				ctx:         context.Background(),
				log:         logger.InitLog(logrus.InfoLevel),
				service:     new(mocks.Service),
				jwtToken:    "token",
				dataKeyWord: "testkey",
				data:        "testdata",
				metadata:    "testmetadata",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Создаем оригинальный os.Stdin и сохраняем его
			originalStdin := os.Stdin

			// Создаем pipe для ввода
			r, w, _ := os.Pipe()

			// Перенаправляем os.Stdin на pipe
			os.Stdin = r
			defer func() {
				os.Stdin = originalStdin
			}()

			// Создаем буфер для вывода данных (для обработки fmt.Println)
			var outputBuf strings.Builder
			wMock := io.MultiWriter(w, io.Writer(&outputBuf))

			go func() {
				_, err := fmt.Fprintln(wMock, tt.args.data)
				assert.NoError(t, err)
				_, err = fmt.Fprintln(wMock, tt.args.dataKeyWord)
				assert.NoError(t, err)
				_, err = fmt.Fprintln(wMock, tt.args.metadata)
				assert.NoError(t, err)
			}()

			data := model.DataBlock{
				DataKeyWord: tt.args.dataKeyWord,
				Data:        tt.args.data,
				MetaData:    tt.args.metadata,
			}
			tt.args.service.On("Add", tt.args.ctx, tt.args.jwtToken, data).Return(nil)

			if err := add(tt.args.ctx, tt.args.log, tt.args.service, tt.args.jwtToken); (err != nil) != tt.wantErr {
				t.Errorf("add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
