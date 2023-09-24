// Пакет service используется в качестве прослойки между
// пакетом хэндлеров и пакетом хранилища
package service

import (
	"context"
	"keeper/internal/utils"

	"github.com/sirupsen/logrus"
)

// Storer - интерфейс взаимодествия с хранилищем
//
//go:generate mockery --name Storer
type Storer interface {
	AddUser(ctx context.Context, login string, password [32]byte) error
	CheckUserAuth(ctx context.Context, login string, password string) error
}

// service - структура, реализующая методы пакета service
type service struct {
	storage Storer
	log     *logrus.Logger
}

func NewService(ctx context.Context, storage Storer,
	log *logrus.Logger) *service {
	return &service{
		storage: storage,
		log:     log,
	}
}

// UserRegister возвращает jwt токен для пользователя, если добавление в бд
// прошло успешно
func (s *service) UserRegister(ctx context.Context, login string,
	password string) (string, error) {

	// добавляем пользователя в бд
	err := s.storage.AddUser(ctx, login, utils.PasswordHash(password))
	if err != nil {
		return "", err
	}

	jwtString, err := utils.GenerateJWTToken(login, s.log)
	if err != nil {
		return "", err
	}
	return jwtString, nil
}

// UserAuthentification проверят логин и пароль пользователя, возвращает jwt токен,
// если все введено верно
func (s *service) UserAuthentification(ctx context.Context, login string,
	password string) (string, error) {

	if err := s.storage.CheckUserAuth(ctx, login, password); err != nil {
		return "", err
	}

	jwtString, err := utils.GenerateJWTToken(login, s.log)
	if err != nil {
		return "", err
	}

	return jwtString, nil
}
