package utils

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"keeper/internal/model"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

// GenerateJWTToken генерирует jwt токен
func GenerateJWTToken(login string, log *logrus.Logger,
	secretPassword string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["login"] = login
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Время истечения токена (1 час)

	jwtString, err := token.SignedString([]byte(secretPassword))
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

// GenerateAESKeyFromPassword генерирует 32-байтовый ключ на основе строки
func GenerateAESKeyFromPassword(password string) ([]byte, error) {

	// Создаем новый хеш SHA-256.
	hash := sha256.New()

	// Записываем пароль в хеш для вычисления хеш-значения.
	_, err := hash.Write([]byte(password))
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

// GCMDataCipher шифрует данные по методу AES-256 GCM
func GCMDataCipher(data string, secretPassword string,
	log *logrus.Logger) ([]byte, error) {

	key, err := GenerateAESKeyFromPassword(secretPassword)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// Создаем AES-256 GCM блок с использованием ключа
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// создаем GCM режим шифрования
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// создаем вектор инициализации из последних байт ключа
	iv := []byte(secretPassword[len(secretPassword)-aesGCM.NonceSize():])

	// шифруем данные
	cipherData := aesGCM.Seal(nil, iv, []byte(data), nil)
	return cipherData, nil
}

// GetLoginFromContext получает логин пользователя из метаданных
func GetLoginFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", model.ErrLoginNotFound
	}
	values := md.Get("login")
	if len(values) == 0 {
		return "", model.ErrLoginNotFound
	}
	login := values[0]
	return login, nil
}
