package handlers

import (
	"context"
	"errors"
	"keeper/internal/model"
	auth "keeper/internal/server/handlers/proto/authService"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
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

// NewHandlerAuth возвращает структуру для операций с хэндлерами
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
	h.log.Debug("Хэндлер для регистрации пользователя")
	jwtString, err := h.service.UserRegister(ctx, in.Login, in.Password)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == pgerrcode.UniqueViolation {
				return nil, status.Errorf(codes.AlreadyExists, model.ErrUserAlreadyExists.Error())
			}
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	response.JwtToken = jwtString

	return &response, nil
}

// UserAuth - хэндлер для аутентификации пользователя
func (h HandlersAuth) UserAuth(ctx context.Context, in *auth.AuthRequest) (
	*auth.AuthResponse, error) {
	var response auth.AuthResponse
	h.log.Debug("Хэндлер для аутентификации пользователя")
	jwtString, err := h.service.UserAuthentification(ctx, in.Login, in.Password)
	if err != nil {
		if errors.Is(err, model.ErrIncorrectPassword) {
			return nil, status.Errorf(codes.Unauthenticated, model.ErrUserAuth.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	response.JwtToken = jwtString
	return &response, nil
}
