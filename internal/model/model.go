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
	ErrLoginNotFound = errors.New("Login not found")
	ErrNotValidToken = errors.New("Not valid token")
)
