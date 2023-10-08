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
	log            *logrus.Logger
	connAuthClient *grpc.ClientConn
	authClient     authservice.AuthServiceClient
	connDataClient *grpc.ClientConn
	dataClient     dataService.DataServiceClient
}

// GetService устанавливает клиентские соединения с gRPC серверами аутентификации и взаимодействия с данными
// и возвращает их в качестве полей структуры service
func GetService(l *logrus.Logger) (*service, error) {
	var service service
	var err error
	service.log = l

	service.connAuthClient, err = grpc.Dial(":9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		l.Error(err.Error())
		return nil, err
	}

	service.authClient = authservice.NewAuthServiceClient(service.connAuthClient)

	service.connDataClient, err = grpc.Dial(":9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		l.Error(err.Error())
		return nil, err
	}

	service.dataClient = dataService.NewDataServiceClient(service.connDataClient)

	return &service, nil
}

// Close закрывает клиентские соединения с gRPC серверами аутентификации и взаимодействия с данными
func (s *service) Close() {
	s.connAuthClient.Close()
	s.connDataClient.Close()
}

// Register передает введенные пользователем логин и пароль в
// в метод регистрации gRPC сервера, получает jwt токен
func (s *service) Register(ctx context.Context, login string, password string) (string, error) {
	requestRegister := &authservice.RegisterRequest{
		Login:    login,
		Password: password,
	}

	resp, err := s.authClient.UserRegister(ctx, requestRegister)
	if err != nil {
		return "", err
	}
	return resp.JwtToken, err
}

// Auth передает введенные пользователем логин и пароль в
// в метод аутентификации gRPC сервера, получает jwt токен
func (s *service) Auth(ctx context.Context, login string, password string) (string, error) {
	requestAuth := &authservice.AuthRequest{
		Login:    login,
		Password: password,
	}

	resp, err := s.authClient.UserAuth(ctx, requestAuth)
	if err != nil {
		return "", err
	}
	return resp.JwtToken, err
}

// Add передает введенные пользователем данные в RPC метод добавления данных
func (s *service) Add(ctx context.Context, jwtToken string, data model.DataBlock) error {

	requestAdd := &dataService.AddingRequest{
		DataKeyWord: data.DataKeyWord,
		Data:        data.Data,
		MetaData:    data.MetaData,
	}

	md := metadata.Pairs("token", jwtToken)
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := s.dataClient.AddData(ctx, requestAdd)
	if err != nil {
		s.log.Error(err.Error())
	}

	return err
}

// Get передает введенный пользователем ключ для идентификации данных в RPC метод получения данных,
// получает данные и метаданные
func (s *service) Get(ctx context.Context, jwtToken string, dataKeyWord string) ([]model.DataBlock, error) {

	requestGet := &dataService.GetRequest{
		DataKeyWord: dataKeyWord,
	}

	md := metadata.Pairs("token", jwtToken)
	ctx = metadata.NewOutgoingContext(ctx, md)

	responseList, err := s.dataClient.GetData(ctx, requestGet)
	if err != nil {
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

// Delete передает введенный пользователем ключ для идентификации данных в RPC метод удаления данных
func (s *service) Delete(ctx context.Context, jwtToken string, dataKeyWord string) error {

	requestDelete := &dataService.DeletionRequest{
		DataKeyWord: dataKeyWord,
	}

	md := metadata.Pairs("token", jwtToken)
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := s.dataClient.DeleteData(ctx, requestDelete)
	if err != nil {
		s.log.Error(err.Error())
	}
	return err
}

// Change передает данные для изменения в RPC метод для изменения данных
func (s *service) Change(ctx context.Context, jwtToken string, data model.DataBlock) error {

	requestChange := &dataService.ChangingRequest{
		DataKeyWord:       data.DataKeyWord,
		DataForChange:     data.Data,
		MetaDataForChange: data.MetaData,
	}

	md := metadata.Pairs("token", jwtToken)
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := s.dataClient.ChangeData(ctx, requestChange)
	if err != nil {
		s.log.Error(err.Error())
	}
	return err
}
