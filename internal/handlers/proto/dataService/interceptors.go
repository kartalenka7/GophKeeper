package dataservice

import (
	"context"
	"errors"
	"keeper/internal/logger"
	"keeper/internal/model"
	"keeper/internal/utils"
	"os"

	"github.com/dgrijalva/jwt-go"
)

// AuthInterceptor функция для интерсептора, которая проверяет jwt токены
func AuthInterceptor(ctx context.Context) (context.Context, error) {
	log := logger.InitLog()
	login, err := utils.GetLoginFromContext(ctx)
	if err != nil {
		log.Error(err.Error())
	}
	var tk model.Token
	goprivate, exist := os.LookupEnv("GOPRIVATE")

	if !exist || goprivate == "" {
		err := errors.New("Заполните переменную среды GOPRIVATE")
		return ctx, err
	}
	token, err := jwt.ParseWithClaims(login, tk, func(t *jwt.Token) (interface{}, error) {
		return []byte(goprivate), nil
	})

	if err != nil {
		log.Error(err.Error())
		return ctx, err
	}
	if !token.Valid {
		return ctx, model.ErrNotValidToken
	}

	return ctx, nil
}
