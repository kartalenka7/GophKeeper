package main

import (
	"context"
	"keeper/internal/config"
	"keeper/internal/logger"
	"keeper/internal/server/handlers"
	"keeper/internal/server/service"
	"keeper/internal/server/storage"
	"net"
	"sync"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	authService "keeper/internal/server/handlers/proto/authService"
	data "keeper/internal/server/handlers/proto/dataService"
)

func main() {

	log := logger.InitLog(logrus.DebugLevel)
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

	var wg sync.WaitGroup

	wg.Add(1)
	// Запускаем сервер авторизации пользователей
	go func() {
		defer wg.Done()
		lisAuth, err := net.Listen("tcp", ":9090")
		if err != nil {
			return
		}

		serverAuth := grpc.NewServer()
		authService.RegisterAuthServiceServer(serverAuth, handlers.NewHandlersAuth(
			service, log))
		log.Info("Запустили gRPC сервис для аутентификации на порте 9090")
		serverAuth.Serve(lisAuth)
	}()

	wg.Add(1)
	// Запускаем сервер взаимодействия с данными
	go func() {
		defer wg.Done()
		lisData, err := net.Listen("tcp", ":9091")
		if err != nil {
			return
		}
		serverData := grpc.NewServer(
			grpc.UnaryInterceptor(auth.UnaryServerInterceptor(data.AuthInterceptor)),
		)

		reflection.Register(serverData)
		data.RegisterDataServiceServer(serverData, handlers.NewHandlersData(service, log))
		log.Info("Запустили gRPC сервис для CRUD операци на порте 9091")
		serverData.Serve(lisData)
	}()

	wg.Wait()
	storage.Close()
}
