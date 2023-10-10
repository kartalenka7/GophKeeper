package handlers

import (
	"context"
	data "keeper/internal/handlers/proto/dataService"
	"keeper/internal/model"
	"keeper/internal/utils"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// HandlersData релизует методы-хэндлеры для CRUD операций с данными
type HandlersData struct {
	data.UnimplementedDataServiceServer
	service Service
	log     *logrus.Logger
}

func NewHandlersData(service Service, log *logrus.Logger) *HandlersData {
	h := &HandlersData{
		service: service,
		log:     log,
	}
	return h
}

// AddData - хэндлер для добавления новых данных для хранения
func (h HandlersData) AddData(ctx context.Context, in *data.AddingRequest) (
	*emptypb.Empty, error) {

	data := model.DataBlock{
		DataKeyWord: in.DataKeyWord,
		DataType:    in.DataType,
		Data:        in.Data,
		MetaData:    in.MetaData,
	}

	login, err := utils.GetLoginFromContext(ctx)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, err.Error())
	}
	data.Login = login

	if err := h.service.AddData(ctx, data); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error in adding data")
	}
	return nil, nil
}

// GetData - хэндлер для получения данных пользователя
func (h HandlersData) GetData(ctx context.Context, in *data.GetRequest) (
	*data.GetResponseList, error) {
	dataResponseList := &data.GetResponseList{}
	var dataGetResponse data.GetResponse

	login, err := utils.GetLoginFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	data, err := h.service.GetData(ctx, login, in.DataKeyWord)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	for _, dataLine := range data {
		dataGetResponse.DataKeyWord = dataLine.DataKeyWord
		dataGetResponse.Data = dataLine.Data
		dataGetResponse.MetaData = dataLine.MetaData
		dataResponseList.Response = append(dataResponseList.Response, &dataGetResponse)
	}
	return dataResponseList, nil
}

// ChangeData - хэндлер для изменения существующих данных пользователя
func (h HandlersData) ChangeData(ctx context.Context, in *data.ChangingRequest) (
	*emptypb.Empty, error) {
	data := model.DataBlock{
		DataKeyWord: in.DataKeyWord,
		Data:        in.DataForChange,
		MetaData:    in.MetaDataForChange,
	}

	login, err := utils.GetLoginFromContext(ctx)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, err.Error())
	}
	data.Login = login
	if err := h.service.ChangeData(ctx, data); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

// DeleteData - хэндлер для удаления данных
func (h HandlersData) DeleteData(ctx context.Context, in *data.DeletionRequest) (
	*emptypb.Empty, error) {

	login, err := utils.GetLoginFromContext(ctx)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, err.Error())
	}

	if err := h.service.DeleteData(ctx, login, in.DataKeyWord); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
