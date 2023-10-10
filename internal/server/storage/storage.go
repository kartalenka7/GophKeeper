package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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

	createDataTable = `CREATE TABLE IF NOT EXISTS dataTable(
						login TEXT,
						dataKeyWord TEXT,
						dataType TEXT,
						data BYTEA,
						metadata TEXT,
						CONSTRAINT fk_login FOREIGN KEY (login) REFERENCES users(login),
       				    CONSTRAINT search_index UNIQUE (login, dataKeyWord)  
						)`
	insertData = `INSERT INTO dataTable(login, dataKeyWord, dataType, data, metadata)
				  VALUES($1, $2, $3, $4, $5)`
	selectData = `SELECT dataKeyWord, dataType, data, metadata 
				  FROM dataTable
				  WHERE login = $1 AND dataKeyWord = $2`
	updateData = `UPDATE dataTable SET data = $1, metadata = $2
				  WHERE login = $3 AND dataKeyWord = $4`
	deleteData = `DELETE FROM dataTable WHERE login = $1 AND dataKeyWord = $2`
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
	log.Debug("Запустили соединение с Postgres, инициализировали таблицы")

	return &storage{
		pgxPool: pool,
		log:     log,
	}, nil
}

// InitTable создает таблицы в бд, если они не существуют
func InitTable(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, createUsersTable); err != nil {
		return err
	}
	_, err := pool.Exec(ctx, createDataTable)
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
	s.log.Debug("Проверяем наличие пользователя в бд")
	var hashPassword []byte

	row := s.pgxPool.QueryRow(ctx, selectPassword, login)
	err := row.Scan(&hashPassword)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}
	inputPasswordHash := sha256.Sum256([]byte(password))

	if hex.EncodeToString(inputPasswordHash[:]) != hex.EncodeToString(hashPassword) {
		err = model.ErrIncorrectPassword
		s.log.Error(err.Error())
		return err
	}
	s.log.Debug("Аутентификация успешна")
	return nil
}

// InsertData добавляет данные пользователя в бд
func (s *storage) InsertData(ctx context.Context, data model.DataBlock) error {
	s.log.Debug("Вставляем строку с данными в таблицу dataTable")
	_, err := s.pgxPool.Exec(ctx, insertData, data.Login, data.DataKeyWord,
		data.DataType, data.CipherData, data.MetaData)
	if err != nil {
		s.log.Error(err.Error())
	}
	return err
}

// GetData выбирает данные пользователя по ключу логин + ключевое слово
func (s *storage) GetData(ctx context.Context, login string,
	dataKeyWord string) ([]model.DataBlock, error) {

	rows, err := s.pgxPool.Query(ctx, selectData, login, dataKeyWord)
	defer rows.Close()
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	var dataBlock model.DataBlock
	var data []model.DataBlock
	for rows.Next() {
		err := rows.Scan(&dataBlock.DataKeyWord, &dataBlock.DataType, &dataBlock.CipherData, &dataBlock.MetaData)
		if err != nil {
			s.log.Error(err.Error())
			return nil, err
		}
		data = append(data, dataBlock)
	}
	if rows.Err() != nil {
		s.log.Error(err.Error())
		return nil, err
	}
	if data == nil {
		err := model.ErrNoRowsSelected
		s.log.Error(err.Error())
		return nil, err
	}
	s.log.Debug("Данные успешно выбраны")
	return data, nil
}

// ChangeData запускает UPDATE на данные пользователя
func (s *storage) ChangeData(ctx context.Context, data model.DataBlock) error {
	_, err := s.pgxPool.Exec(ctx, updateData, data.CipherData, data.MetaData,
		data.Login, data.DataKeyWord)
	if err != nil {
		s.log.Error(err.Error())
	}
	return err
}

// DeleteData удаляет данные пользователя по ключу
func (s *storage) DeleteData(ctx context.Context, login string, dataKeyWord string) error {
	_, err := s.pgxPool.Exec(ctx, deleteData, login, dataKeyWord)
	if err != nil {
		s.log.Error(err.Error())
	}
	return err
}

func (s *storage) Close() {
	s.pgxPool.Close()
}
