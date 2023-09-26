package main

import (
	"context"
	"keeper/internal/config"
	"keeper/internal/handlers"
	"keeper/internal/logger"
	"keeper/internal/service"
	"keeper/internal/storage"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	authService "keeper/internal/handlers/proto/authService"
	data "keeper/internal/handlers/proto/dataService"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
)

func main() {

	log := logger.InitLog()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := config.GetConfig(log)
	if err != nil {
		return
	}

	storage, err := storage.NewStorage(ctx, log, config)
	if err != nil {
		return
	}

	service := service.NewService(ctx, storage, log, config)

	// Запускаем сервер авторизации пользователей
	lisAuth, err := net.Listen("tcp", ":9090")
	if err != nil {
		return
	}
	serverAuth := grpc.NewServer()
	authService.RegisterAuthServiceServer(serverAuth, handlers.NewHandlersAuth(
		service, log))
	serverAuth.Serve(lisAuth)

	// Запускаем сервер взаимодействия с данными
	lisData, err := net.Listen("tcp", ":9091")
	if err != nil {
		return
	}
	serverData := grpc.NewServer(
		grpc.UnaryInterceptor(auth.UnaryServerInterceptor(data.AuthInterceptor)),
	)
	reflection.Register(serverData)
	data.RegisterDataServiceServer(serverData, handlers.NewHandlersData(service, log))
	serverAuth.Serve(lisData)
}
