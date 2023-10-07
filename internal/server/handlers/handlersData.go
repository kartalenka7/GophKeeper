package handlers

import (
	"context"
	"errors"
	"keeper/internal/model"
	data "keeper/internal/server/handlers/proto/dataService"

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
	h.log.Debug("Хэндлер для добавления данных")

	data := model.DataBlock{
		DataKeyWord: in.DataKeyWord,
		DataType:    in.DataType,
		Data:        in.Data,
		MetaData:    in.MetaData,
	}

	if err := h.service.AddData(ctx, data); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error in adding data")
	}
	return &emptypb.Empty{}, nil
}

// GetData - хэндлер для получения данных пользователя
func (h HandlersData) GetData(ctx context.Context, in *data.GetRequest) (
	*data.GetResponseList, error) {
	h.log.Debug("Хэндлер для получения данных")
	dataResponseList := &data.GetResponseList{}
	var dataGetResponse data.GetResponse

	data, err := h.service.GetData(ctx, in.DataKeyWord)
	if err != nil {
		if errors.Is(err, model.ErrNoRowsSelected) {
			return nil, status.Error(codes.NotFound, model.ErrNoRowsSelected.Error())
		}
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
	h.log.Debug("Хэндлер для изменения данных")
	data := model.DataBlock{
		DataKeyWord: in.DataKeyWord,
		Data:        in.DataForChange,
		MetaData:    in.MetaDataForChange,
	}

	if err := h.service.ChangeData(ctx, data); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

// DeleteData - хэндлер для удаления данных
func (h HandlersData) DeleteData(ctx context.Context, in *data.DeletionRequest) (
	*emptypb.Empty, error) {
	h.log.Debug("Хэндлер для удаления данных")

	if err := h.service.DeleteData(ctx, in.DataKeyWord); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
