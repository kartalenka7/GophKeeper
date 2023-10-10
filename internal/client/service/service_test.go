package service

import (
	"context"
	"keeper/internal/client/service/mocks"
	"keeper/internal/logger"
	"keeper/internal/model"
	authservice "keeper/internal/server/handlers/proto/authService"
	dataService "keeper/internal/server/handlers/proto/dataService"
	dataservice "keeper/internal/server/handlers/proto/dataService"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestClientServiceRegister(t *testing.T) {

	mockServiceClient := new(mocks.AuthServiceClient)

	type args struct {
		ctx      context.Context
		login    string
		password string
	}
	tests := []struct {
		name    string
		s       *service
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Успешная отправка запроса на регистрацию",
			s: &service{
				log:        logger.InitLog(logrus.InfoLevel),
				authClient: mockServiceClient,
			},
			args: args{
				ctx:      context.Background(),
				login:    "user",
				password: "password",
			},
			wantErr: false,
		},
	}
	response := &authservice.RegisterResponse{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			requestRegister := &authservice.RegisterRequest{
				Login:    tt.args.login,
				Password: tt.args.password,
			}
			response.JwtToken = "token"
			mockServiceClient.On("UserRegister", tt.args.ctx, requestRegister).Return(response, nil)

			_, err := tt.s.Register(tt.args.ctx, tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClientServiceAuth(t *testing.T) {

	mockServiceClient := new(mocks.AuthServiceClient)

	type args struct {
		ctx      context.Context
		login    string
		password string
	}
	tests := []struct {
		name    string
		s       *service
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Успешная отправка запроса на аутентификацию",
			s: &service{
				log:        logger.InitLog(logrus.InfoLevel),
				authClient: mockServiceClient,
			},
			args: args{
				ctx:      context.Background(),
				login:    "user",
				password: "password",
			},
			wantErr: false,
		},
	}
	response := &authservice.AuthResponse{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestAuth := &authservice.AuthRequest{
				Login:    tt.args.login,
				Password: tt.args.password,
			}
			response.JwtToken = "token"
			mockServiceClient.On("UserAuth", tt.args.ctx, requestAuth).Return(response, nil)

			_, err := tt.s.Auth(tt.args.ctx, tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Auth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClientServiceAdd(t *testing.T) {

	mockServiceClient := new(mocks.DataServiceClient)

	type args struct {
		ctx      context.Context
		jwtToken string
		data     model.DataBlock
	}
	tests := []struct {
		name    string
		s       *service
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Успешная отправка запроса на добавление данных",
			s: &service{
				log:        logger.InitLog(logrus.InfoLevel),
				dataClient: mockServiceClient,
			},
			args: args{
				ctx:      context.Background(),
				jwtToken: "token",
				data: model.DataBlock{
					DataKeyWord: "key",
					Data:        "data",
					MetaData:    "metadata",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			requestAdd := &dataService.AddingRequest{
				DataKeyWord: tt.args.data.DataKeyWord,
				Data:        tt.args.data.Data,
				MetaData:    tt.args.data.MetaData,
			}

			md := metadata.Pairs("token", tt.args.jwtToken)
			ctx := metadata.NewOutgoingContext(tt.args.ctx, md)
			mockServiceClient.On("AddData", ctx, requestAdd).Return(&emptypb.Empty{}, nil)

			if err := tt.s.Add(tt.args.ctx, tt.args.jwtToken, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("service.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientServiceGet(t *testing.T) {
	mockServiceClient := new(mocks.DataServiceClient)
	type args struct {
		ctx         context.Context
		jwtToken    string
		dataKeyWord string
	}
	tests := []struct {
		name    string
		s       *service
		args    args
		want    []model.DataBlock
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Успешная отправка запроса на получение данных",
			s: &service{
				log:        logger.InitLog(logrus.InfoLevel),
				dataClient: mockServiceClient,
			},
			args: args{
				ctx:         context.Background(),
				jwtToken:    "token",
				dataKeyWord: "key",
			},
			want: []model.DataBlock{
				{
					DataKeyWord: "key",
					Data:        "data",
					MetaData:    "metadata",
				},
			},
			wantErr: false,
		},
	}
	responseList := &dataservice.GetResponseList{
		Response: []*dataservice.GetResponse{
			{
				DataKeyWord: "key",
				Data:        "data",
				MetaData:    "metadata",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			md := metadata.Pairs("token", tt.args.jwtToken)
			ctx := metadata.NewOutgoingContext(tt.args.ctx, md)

			requestGet := &dataService.GetRequest{
				DataKeyWord: tt.args.dataKeyWord,
			}
			mockServiceClient.On("GetData", ctx, requestGet).Return(responseList, nil)

			got, err := tt.s.Get(tt.args.ctx, tt.args.jwtToken, tt.args.dataKeyWord)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}

func TestClientServiceDelete(t *testing.T) {
	mockServiceClient := new(mocks.DataServiceClient)
	type args struct {
		ctx         context.Context
		jwtToken    string
		dataKeyWord string
	}
	tests := []struct {
		name    string
		s       *service
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Успешная отправка запроса на удаление данных",
			s: &service{
				log:        logger.InitLog(logrus.InfoLevel),
				dataClient: mockServiceClient,
			},
			args: args{
				ctx:         context.Background(),
				jwtToken:    "token",
				dataKeyWord: "key",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			md := metadata.Pairs("token", tt.args.jwtToken)
			ctx := metadata.NewOutgoingContext(tt.args.ctx, md)

			requestDelete := &dataService.DeletionRequest{
				DataKeyWord: tt.args.dataKeyWord,
			}
			mockServiceClient.On("DeleteData", ctx, requestDelete).Return(&emptypb.Empty{}, nil)

			if err := tt.s.Delete(tt.args.ctx, tt.args.jwtToken, tt.args.dataKeyWord); (err != nil) != tt.wantErr {
				t.Errorf("service.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientServiceChange(t *testing.T) {
	mockServiceClient := new(mocks.DataServiceClient)
	type args struct {
		ctx      context.Context
		jwtToken string
		data     model.DataBlock
	}
	tests := []struct {
		name    string
		s       *service
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Успешная отправка запроса на изменение",
			s: &service{
				log:        logger.InitLog(logrus.InfoLevel),
				dataClient: mockServiceClient,
			},
			args: args{
				ctx:      context.Background(),
				jwtToken: "token",
				data: model.DataBlock{
					DataKeyWord: "key",
					Data:        "data",
					MetaData:    "metadata",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			md := metadata.Pairs("token", tt.args.jwtToken)
			ctx := metadata.NewOutgoingContext(tt.args.ctx, md)

			requestChange := &dataService.ChangingRequest{
				DataKeyWord:       tt.args.data.DataKeyWord,
				DataForChange:     tt.args.data.Data,
				MetaDataForChange: tt.args.data.MetaData,
			}
			mockServiceClient.On("ChangeData", ctx, requestChange).Return(&emptypb.Empty{}, nil)

			if err := tt.s.Change(tt.args.ctx, tt.args.jwtToken, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("service.Change() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetService(t *testing.T) {
	type args struct {
		l *logrus.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Получение объекта структуры сервис",
			args: args{
				l: logger.InitLog(logrus.InfoLevel),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetService(tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}
