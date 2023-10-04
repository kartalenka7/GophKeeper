package utils

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"keeper/internal/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

// GenerateJWTToken генерирует jwt токен
func GenerateJWTToken(login string, log *logrus.Logger,
	secretPassword string) (string, error) {
	log.Debug("Генерируем JWT токен")

	tk := &model.Token{Login: login}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
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

// prepareAESGCM подготавливает режим AES-256 GCM,
// возвращает вектор инициализации
func prepareAESGCM(log *logrus.Logger,
	secretPassword string) (cipher.AEAD, []byte, error) {

	// Создаем новый хеш SHA-256.
	hash := sha256.New()

	// Записываем пароль в хеш для вычисления хеш-значения.
	_, err := hash.Write([]byte(secretPassword))
	if err != nil {
		log.Error(err.Error())
		return nil, nil, err
	}

	key := hash.Sum(nil)

	// Создаем AES-256 GCM блок с использованием ключа
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error(err.Error())
		return nil, nil, err
	}

	// создаем GCM режим шифрования
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Error(err.Error())
		return nil, nil, err
	}

	// создаем вектор инициализации из последних байт ключа
	iv := []byte(secretPassword[len(secretPassword)-aesGCM.NonceSize():])
	return aesGCM, iv, nil
}

// GCMDataCipher шифрует данные по методу AES-256 GCM
func GCMDataCipher(data string, secretPassword string,
	log *logrus.Logger) ([]byte, error) {

	aesGCM, nonce, err := prepareAESGCM(log, secretPassword)
	if err != nil {
		return nil, err
	}
	// шифруем данные
	cipherData := aesGCM.Seal(nil, nonce, []byte(data), nil)
	return cipherData, nil
}

// GCMDataDecipher дешифрует данные по методу AES-256 GCM
func GCMDataDecipher(cipherData []byte, secretPassword string, log *logrus.Logger) (string, error) {
	log.Debug("Дешифруем данные")

	aesGCM, iv, err := prepareAESGCM(log, secretPassword)
	if err != nil {
		return "", err
	}

	// Расшифровать данные
	plainData, err := aesGCM.Open(nil, iv, cipherData, nil)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	log.WithFields(logrus.Fields{
		"plaintext": string(plainData),
	}).Debug("Дешифрованные данные")
	return string(plainData), nil
}

// GetLoginFromContext получает логин пользователя из метаданных контекста
func GetLoginFromContext(ctx context.Context, secretPassword string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", model.ErrTokenNotFound
	}
	values := md.Get("token")
	if len(values) == 0 {
		return "", model.ErrTokenNotFound
	}
	jwtString := values[0]

	tk := model.Token{}
	token, err := jwt.ParseWithClaims(jwtString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretPassword), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", model.ErrNotValidToken
	}
	return tk.Login, nil
}
