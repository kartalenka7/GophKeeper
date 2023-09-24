package handlers

import (
	"context"
	data "keeper/internal/handlers/proto/dataService"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

// HandlersData релизует методы-хэндлеры для CRUD операций с данными
type HandlersData struct {
	data.UnimplementedDataServiceServer
	service Service
	log     *logrus.Logger
}

// NewRouter инициализирует роутер chi и маршрутизирует
// все запросы в подходящие методы-хэндлеры
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
	return nil, nil
}

// ChangeData - хэндлер для изменения существующих данных пользователя
func (h HandlersData) ChangeData(ctx context.Context, in *data.ChangingRequest) (
	*emptypb.Empty, error) {
	return nil, nil
}

// GetData - хэндлер для получения данных пользователя
func (h HandlersData) GetData(ctx context.Context, in *data.GetRequest) (
	*data.GetResponse, error) {
	return nil, nil
}

// DeleteData - хэндлер для удаления данных
func (h HandlersData) DeleteData(ctx context.Context, in *data.DeletionRequest) (
	*emptypb.Empty, error) {
	return nil, nil
}
