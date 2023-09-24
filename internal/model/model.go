package model

import "github.com/dgrijalva/jwt-go"

type Token struct {
	Login string
	jwt.StandardClaims
}

type Config struct {
	ConfigFile string `env:"CONFIG"`
	Database   string `json:"database_conn"`
}
