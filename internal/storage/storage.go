package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"keeper/internal/model"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type storage struct {
	pgxPool *pgxpool.Pool
	log     *logrus.Logger
}

var (
	createUsersTable = `CREATE TABLE IF NOT EXISTS users(
						login TEXT PRIMARY KEY,
						password BYTEA
					    );`
	insertUser     = `INSERT INTO users(login, password) VALUES($1, $2)`
	selectPassword = `SELECT password FROM users WHERE login = $1`
)

// NewStorage инициализирует пул соединений с базой данных
func NewStorage(ctx context.Context, log *logrus.Logger,
	config model.Config) (
	*storage, error) {

	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pool, err := pgxpool.Connect(ctxTimeout, config.Database)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	err = InitTable(ctxTimeout, pool)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &storage{
		pgxPool: pool,
		log:     log,
	}, nil
}

// InitTable создает таблицы в бд, если они не существуют
func InitTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, createUsersTable)
	return err
}

// AddUser добавляет нового пользователя в бд
func (s *storage) AddUser(ctx context.Context, login string,
	password [32]byte) error {

	// Convert the byte slice to hexadecimal format
	hexEncodedPassword := fmt.Sprintf("\\x%x", password)

	_, err := s.pgxPool.Exec(ctx, insertUser, login, hexEncodedPassword)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}
	return nil
}

// CheckUserAuth проверяет логин и хэш пароля пользователя по таблице users
func (s *storage) CheckUserAuth(ctx context.Context, login string,
	password string) error {
	var hashPassword []byte

	row := s.pgxPool.QueryRow(ctx, selectPassword, login)
	err := row.Scan(&hashPassword)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}
	inputPasswordHash := sha256.Sum256([]byte(password))

	if hex.EncodeToString(inputPasswordHash[:]) != hex.EncodeToString(hashPassword) {
		err = errors.New(`incorrect password`)
		s.log.Error(err.Error())
		return err
	}
	return nil
}
