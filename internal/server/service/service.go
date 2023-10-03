// Пакет service используется в качестве прослойки между
// пакетом хэндлеров и пакетом хранилища
package service

import (
	"context"
	"keeper/internal/model"
	"keeper/internal/utils"

	"github.com/sirupsen/logrus"
)

// Storer - интерфейс взаимодествия с хранилищем
//
//go:generate mockery --name Storer
type Storer interface {
	AddUser(ctx context.Context, login string, password [32]byte) error
	CheckUserAuth(ctx context.Context, login string, password string) error
	InsertData(ctx context.Context, data model.DataBlock) error
	GetData(ctx context.Context, login string, dataKeyWord string) ([]model.DataBlock, error)
	ChangeData(ctx context.Context, data model.DataBlock) error
	DeleteData(ctx context.Context, login string, dataKeyWord string) error
}

// service - структура, реализующая методы пакета service
type service struct {
	storage Storer
	log     *logrus.Logger
	config  model.Config
}

func NewService(ctx context.Context, storage Storer,
	log *logrus.Logger, cfg model.Config) *service {
	return &service{
		storage: storage,
		log:     log,
		config:  cfg,
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

	jwtString, err := utils.GenerateJWTToken(login, s.log, s.config.SecretPassword)
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

	jwtString, err := utils.GenerateJWTToken(login, s.log, s.config.SecretPassword)
	if err != nil {
		return "", err
	}

	return jwtString, nil
}

// AddData шифрует данные и отправляет их в storage
func (s *service) AddData(ctx context.Context, data model.DataBlock) error {
	cipherData, err := utils.GCMDataCipher(data.Data, s.config.SecretPassword, s.log)
	if err != nil {
		return err
	}
	data.CipherData = cipherData

	login, err := utils.GetLoginFromContext(ctx, s.config.SecretPassword)
	if err != nil {
		return err
	}

	data.Login = login
	err = s.storage.InsertData(ctx, data)
	return err
}

// GetData возвращает данные пользователя
func (s *service) GetData(ctx context.Context,
	dataKeyWord string) ([]model.DataBlock, error) {

	login, err := utils.GetLoginFromContext(ctx, s.config.SecretPassword)
	if err != nil {
		return nil, err
	}
	data, err := s.storage.GetData(ctx, login, dataKeyWord)
	if err != nil {
		return nil, err
	}
	var dataReturn []model.DataBlock
	for _, dataLine := range data {
		dataDecipher, err := utils.GCMDataDecipher(dataLine.CipherData, s.config.SecretPassword,
			s.log)
		if err != nil {
			return nil, err
		}

		dataReturn = append(dataReturn, model.DataBlock{
			DataKeyWord: dataLine.DataKeyWord,
			Data:        dataDecipher,
			MetaData:    dataLine.MetaData,
		})
	}
	return dataReturn, err
}

// ChangeData шифрует новые данные и отправляет их в storage
func (s *service) ChangeData(ctx context.Context, dataForChange model.DataBlock) error {
	cipherData, err := utils.GCMDataCipher(dataForChange.Data, s.config.SecretPassword, s.log)
	if err != nil {
		return err
	}
	dataForChange.CipherData = cipherData

	login, err := utils.GetLoginFromContext(ctx, s.config.SecretPassword)
	if err != nil {
		return err
	}

	dataForChange.Login = login
	return s.storage.ChangeData(ctx, dataForChange)
}

func (s *service) DeleteData(ctx context.Context, dataKeyWord string) error {

	login, err := utils.GetLoginFromContext(ctx, s.config.SecretPassword)
	if err != nil {
		return err
	}

	return s.storage.DeleteData(ctx, login, dataKeyWord)
}
