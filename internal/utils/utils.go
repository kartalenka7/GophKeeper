package utils

import (
	"crypto/sha256"
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

// GenerateJWTToken генерирует jwt токен
func GenerateJWTToken(login string, log *logrus.Logger) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["login"] = login
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Время истечения токена (1 час)

	goprivate, exist := os.LookupEnv("GOPRIVATE")

	if !exist || goprivate == "" {
		err := errors.New("Заполните переменную среды GOPRIVATE")
		return "", err
	}

	jwtString, err := token.SignedString([]byte(goprivate))
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	return jwtString, nil
}

// PasswordHash возвращает хэш пароля по методу SHA-256
func PasswordHash(password string) [32]byte {
	return sha256.Sum256([]byte(password))
}
