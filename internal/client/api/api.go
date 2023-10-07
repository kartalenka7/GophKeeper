package api

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"keeper/internal/model"
)

type Service interface {
	Register(ctx context.Context, login string, password string) (string, error)
	Auth(ctx context.Context, login string, password string) (string, error)
	Add(ctx context.Context, jwtToken string, data model.DataBlock) error
	Get(ctx context.Context, jwtToken string, dataKeyWord string) ([]model.DataBlock, error)
	Delete(ctx context.Context, jwtToken string, dataKeyWord string) error
	Change(ctx context.Context, jwtToken string, data model.DataBlock) error
	/*checkData() // проверить размер файлов */
}

func InitCLIApp(ctx context.Context, log *logrus.Logger, service Service) *cli.App {
	app := cli.NewApp()
	app.Name = "Веб приложение для хранения паролей"

	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "Запустить приложение",
			Action: func(c *cli.Context) error {

				fmt.Println("Приложение запущено. Для выхода введите exit")
				var jwtToken string
				var err error
				for {
					var input string
					fmt.Scanln(&input)
					switch input {
					case "exit":
						return nil
					case "register":
						jwtToken, err = register(ctx, log, service)
						if err != nil {
							return err
						}
					case "auth":
						jwtToken, err = auth(ctx, log, service)
						if err != nil {
							return err
						}
					case "add":
						if checkAuth(jwtToken, log) {
							continue
						}
						if err = add(ctx, log, service, jwtToken); err != nil {
							return err
						}
					case "get":
						if checkAuth(jwtToken, log) {
							continue
						}
						if err = get(ctx, log, service, jwtToken); err != nil {
							return err
						}
					case "delete":
						if checkAuth(jwtToken, log) {
							continue
						}
						if err = delete(ctx, log, service, jwtToken); err != nil {
							return err
						}
					case "change":
						if checkAuth(jwtToken, log) {
							continue
						}
						if err = change(ctx, log, service, jwtToken); err != nil {
							return err
						}
					default:
						fmt.Println("register - регистрация пользователя")
						fmt.Println("auth - аутентификация пользователя")
						fmt.Println("add - добавить данные")
						fmt.Println("get - получить данные")
						fmt.Println("change - изменить данные")
						fmt.Println("delete - удалить данные")
					}
				}
			},
		},
	}
	return app
}

func checkAuth(jwtToken string, log *logrus.Logger) bool {
	if jwtToken != "" {
		return false
	}
	fmt.Println(model.ErrNoAuthentification.Error())
	return true
}

func register(ctx context.Context, log *logrus.Logger,
	service Service) (string, error) {
	var login string
	var password string
	fmt.Println("Введите логин")
	_, err := fmt.Scanln(&login)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	fmt.Println("Введите пароль")
	_, err = fmt.Scanln(&password)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	jwtToken, err := service.Register(ctx, login, password)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.AlreadyExists:
				fmt.Println(e.Message())
				return "", nil
			default:
				log.Error(err.Error())
				return "", err
			}
		}
		log.Error(err.Error())
		return "", err
	}
	fmt.Println("Регистрация успешна")
	return jwtToken, nil
}

func auth(ctx context.Context, log *logrus.Logger,
	service Service) (string, error) {
	var login string
	var password string
	fmt.Println("Введите логин")
	_, err := fmt.Scanln(&login)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	fmt.Println("Введите пароль")
	_, err = fmt.Scanln(&password)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	jwtToken, err := service.Auth(ctx, login, password)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Unauthenticated:
				fmt.Println(e.Message())
				return "", nil
			default:
				log.Error(err.Error())
				return "", err
			}
		}
		log.Error(err.Error())
		return "", err
	}
	fmt.Println("Аутентификация успешна")
	return jwtToken, nil
}

func add(ctx context.Context, log *logrus.Logger,
	service Service, jwtToken string) error {
	fmt.Println("Введите строку с данными для добавления или путь к файлу")
	var data model.DataBlock
	_, err := fmt.Scanln(&data.Data)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if strings.Contains(data.Data, `\`) {
		data.Data, err = readFromFile(data.Data, log)
	}
	fmt.Println("Введите ключ для однозначной идентификации данных")
	_, err = fmt.Scanln(&data.DataKeyWord)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	fmt.Println(`Введите дополнительные метаданные (не рекомендуется вводить чувствительную информацию),` +
		`если необходимо`)
	_, err = fmt.Scan(&data.MetaData)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	fmt.Printf("Вы ввели %s\n", data)
	err = service.Add(ctx, jwtToken, data)
	if err != nil {
		return err
	}
	fmt.Println("Данные успешно добавлены")
	return nil
}

func readFromFile(fileName string, log *logrus.Logger) (string, error) {

	// Открываем файл для чтения.
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		log.Fatalf("Ошибка при открытии файла: %v", err)
	}
	defer file.Close()

	// Определяем размер файла.
	fileInfo, err := file.Stat()
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	fileSize := fileInfo.Size()
	if fileSize > 3072 {
		err = model.ErrBigFile
		log.Error(err.Error())
		return "", err
	}

	// Создаем буфер для чтения данных из файла.
	data := make([]byte, fileSize)

	// Читаем данные из файла в буфер.
	_, err = file.Read(data)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	return string(data), err
}

func get(ctx context.Context, log *logrus.Logger,
	service Service, jwtToken string) error {
	var keyWord string
	fmt.Println("Введите ключ для однозначной идентификации данных")
	_, err := fmt.Scanln(&keyWord)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	data, err := service.Get(ctx, jwtToken, keyWord)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				fmt.Println(e.Message())
				return nil
			default:
				log.Error(err.Error())
				return err
			}
		}
		log.Error(err.Error())
		return err
	}
	fmt.Printf("Сохраненные данные: %s\n", data)
	return nil
}

func delete(ctx context.Context, log *logrus.Logger,
	service Service, jwtToken string) error {
	var keyWord string
	fmt.Println("Введите ключ для однозначной идентификации данных")
	_, err := fmt.Scanln(&keyWord)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	err = service.Delete(ctx, jwtToken, keyWord)
	if err != nil {
		return err
	}
	fmt.Println("Запись удалена")
	return nil
}

func change(ctx context.Context, log *logrus.Logger,
	service Service, jwtToken string) error {
	var data model.DataBlock
	fmt.Println("Введите ключ для однозначной идентификации данных")
	_, err := fmt.Scanln(&data.DataKeyWord)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	fmt.Println("Введите данные для изменения")
	_, err = fmt.Scanln(&data.Data)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	fmt.Println(`Введите дополнительные метаданные (не рекомендуется вводить чувствительную информацию),` +
		`если необходимо`)
	_, err = fmt.Scan(&data.MetaData)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	fmt.Printf("Вы ввели %s", data)

	err = service.Change(ctx, jwtToken, data)
	if err != nil {
		return err
	}
	fmt.Println("Данные успешно изменены")
	return nil
}
