package service

import (
	"context"
	"keeper/internal/model"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	authservice "keeper/internal/server/handlers/proto/authService"
	dataService "keeper/internal/server/handlers/proto/dataService"
)

type service struct {
	log *logrus.Logger
}

func GetService(l *logrus.Logger) *service {
	return &service{
		log: l,
	}
}

func (s *service) Register(ctx context.Context, login string, password string) (string, error) {
	conn, err := grpc.Dial(":9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		s.log.Error(err.Error())
		return "", err
	}

	client := authservice.NewAuthServiceClient(conn)

	requestRegister := &authservice.RegisterRequest{
		Login:    login,
		Password: password,
	}

	// в response добавить jwt токен
	resp, err := client.UserRegister(ctx, requestRegister)
	if err != nil {
		s.log.Error(err.Error())
	}
	return resp.JwtToken, err
}

func (s *service) Auth(ctx context.Context, login string, password string) (string, error) {
	conn, err := grpc.Dial(":9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		s.log.Error(err.Error())
		return "", err
	}

	client := authservice.NewAuthServiceClient(conn)

	requestAuth := &authservice.AuthRequest{
		Login:    login,
		Password: password,
	}

	// в response добавить jwt токен
	resp, err := client.UserAuth(ctx, requestAuth)
	if err != nil {
		s.log.Error(err.Error())
	}
	return resp.JwtToken, err
}

func (s *service) Add(ctx context.Context, jwtToken string, data model.DataBlock) error {
	conn, err := grpc.Dial(":9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	client := dataService.NewDataServiceClient(conn)

	requestAdd := &dataService.AddingRequest{
		DataKeyWord: data.DataKeyWord,
		Data:        data.Data,
		MetaData:    data.MetaData,
	}

	md := metadata.Pairs("token", jwtToken)
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err = client.AddData(ctx, requestAdd)
	if err != nil {
		s.log.Error(err.Error())
	}

	return err
}

func (s *service) Get(ctx context.Context, jwtToken string, dataKeyWord string) ([]model.DataBlock, error) {
	conn, err := grpc.Dial(":9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}
	client := dataService.NewDataServiceClient(conn)

	requestGet := &dataService.GetRequest{
		DataKeyWord: dataKeyWord,
	}

	md := metadata.Pairs("token", jwtToken)
	ctx = metadata.NewOutgoingContext(ctx, md)

	responseList, err := client.GetData(ctx, requestGet)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	var data []model.DataBlock

	for _, resp := range responseList.Response {
		data = append(data, model.DataBlock{
			DataKeyWord: resp.DataKeyWord,
			Data:        resp.Data,
			MetaData:    resp.MetaData,
		})
	}
	return data, nil
}

func (s *service) Delete(ctx context.Context, jwtToken string, dataKeyWord string) error {
	conn, err := grpc.Dial(":9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		s.log.Error(err.Error())
		return err
	}
	client := dataService.NewDataServiceClient(conn)

	requestDelete := &dataService.DeletionRequest{
		DataKeyWord: dataKeyWord,
	}

	md := metadata.Pairs("token", jwtToken)
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err = client.DeleteData(ctx, requestDelete)
	if err != nil {
		s.log.Error(err.Error())
	}
	return err
}

func (s *service) Change(ctx context.Context, jwtToken string, data model.DataBlock) error {
	conn, err := grpc.Dial(":9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		s.log.Error(err.Error())
		return err
	}
	client := dataService.NewDataServiceClient(conn)

	requestChange := &dataService.ChangingRequest{
		DataKeyWord:       data.DataKeyWord,
		DataForChange:     data.Data,
		MetaDataForChange: data.MetaData,
	}

	md := metadata.Pairs("token", jwtToken)
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err = client.ChangeData(ctx, requestChange)
	if err != nil {
		s.log.Error(err.Error())
	}
	return err
}
