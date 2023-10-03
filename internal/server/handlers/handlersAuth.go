package handlers

import (
	"context"
	"keeper/internal/model"
	auth "keeper/internal/server/handlers/proto/authService"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	UserRegister(ctx context.Context, login string, password string) (string, error)
	UserAuthentification(ctx context.Context, login string, password string) (string, error)
	AddData(ctx context.Context, data model.DataBlock) error
	GetData(ctx context.Context, dataKeyWord string) ([]model.DataBlock, error)
	ChangeData(ctx context.Context, dataForChange model.DataBlock) error
	DeleteData(ctx context.Context, dataKeyWord string) error
}

// HandlerAuth реализует методы-хэндлеры регистрации
// и аутентификации пользователя
type HandlersAuth struct {
	auth.UnimplementedAuthServiceServer
	service Service
	log     *logrus.Logger
}

// NewHandlerAuth
func NewHandlersAuth(service Service, log *logrus.Logger) *HandlersAuth {
	h := &HandlersAuth{
		service: service,
		log:     log,
	}
	return h
}

// UserRegister - хэндлер для регистрации пользователя
func (h HandlersAuth) UserRegister(ctx context.Context, in *auth.RegisterRequest) (
	*auth.RegisterResponse, error) {
	var response auth.RegisterResponse

	jwtString, err := h.service.UserRegister(ctx, in.Login, in.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error in user registration")
	}

	response.JwtToken = jwtString

	return &response, nil
}

// UserAuth - хэндлер для аутентификации пользователя
func (h HandlersAuth) UserAuth(ctx context.Context, in *auth.AuthRequest) (
	*auth.AuthResponse, error) {
	var response auth.AuthResponse
	h.log.Info("Хэндлер для аутентификации пользователя")
	jwtString, err := h.service.UserAuthentification(ctx, in.Login, in.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error in user authentification")
	}
	response.JwtToken = jwtString
	return &response, nil
}
