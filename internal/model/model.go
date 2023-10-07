package model

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	Login string
	jwt.StandardClaims
}

// Config - структура с данными конфигурации приложения
type Config struct {
	ConfigFile     string `env:"CONFIG"`
	Database       string `json:"database_conn"`
	SecretPassword string
}

// DataBlock - структура для операций с данными пользователя
type DataBlock struct {
	Login       string
	DataKeyWord string
	DataType    string
	Data        string
	CipherData  []byte
	MetaData    string
}

var (
	ErrLoginNotFound      = errors.New("Login not found")
	ErrTokenNotFound      = errors.New("Token not found")
	ErrNotValidToken      = errors.New("Not valid token")
	ErrNoRowsSelected     = errors.New("No rows selected")
	ErrUniqueViolation    = errors.New("Login unique violation")
	ErrUserNotFound       = errors.New("User not found")
	ErrUserRegister       = errors.New("error in user registration")
	ErrUserAlreadyExists  = errors.New("Пользователь уже существует, придумайте другой логин")
	ErrUserAuth           = errors.New("Неправильный логин или пароль")
	ErrNoAuthentification = errors.New("Сначала пройдите регистрацию или аутентификацию")
	ErrBigFile            = errors.New("Слишком большой файл")
	ErrIncorrectPassword  = errors.New("incorrect login or password")
)
