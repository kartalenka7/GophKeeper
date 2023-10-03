package dataservice

import (
	"context"
	"errors"
	"keeper/internal/logger"
	"keeper/internal/utils"
	"os"

	"github.com/sirupsen/logrus"
)

// AuthInterceptor функция для интерсептора, которая проверяет jwt токены
func AuthInterceptor(ctx context.Context) (context.Context, error) {
	log := logger.InitLog(logrus.DebugLevel)
	log.Info("Интерсептор с проверкой jwt токена")

	goprivate, exist := os.LookupEnv("GOPRIVATE")

	if !exist || goprivate == "" {
		err := errors.New("Заполните переменную среды GOPRIVATE")
		return ctx, err
	}

	_, err := utils.GetLoginFromContext(ctx, goprivate)
	if err != nil {
		log.Error(err.Error())
		return ctx, err
	}

	return ctx, nil
}
