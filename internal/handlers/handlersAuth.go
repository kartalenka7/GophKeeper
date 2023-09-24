package handlers

import (
	"context"
	auth "keeper/internal/handlers/proto/authService"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service interface {
	UserRegister(ctx context.Context, login string, password string) (string, error)
	UserAuthentification(ctx context.Context, login string, password string) (string, error)
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
	*emptypb.Empty, error) {

	jwtString, err := h.service.UserRegister(ctx, in.Login, in.Password)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error in user registration")
	}

	md := metadata.Pairs("authorization", jwtString)
	ctx = metadata.NewOutgoingContext(ctx, md)

	return &emptypb.Empty{}, nil
}

// UserAuth - хэндлер для аутентификации пользователя
func (h HandlersAuth) UserAuth(ctx context.Context, in *auth.AuthRequest) (
	*emptypb.Empty, error) {

	jwtString, err := h.service.UserAuthentification(ctx, in.Login, in.Password)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error in user authentification")
	}
	
	md := metadata.Pairs("authorization", jwtString)
	ctx = metadata.NewOutgoingContext(ctx, md)

	return &emptypb.Empty{}, nil
}
