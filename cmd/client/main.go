package main

import (
	"context"
	"keeper/internal/client/api"
	"keeper/internal/client/service"
	"keeper/internal/logger"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logger.InitLog(logrus.DebugLevel)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := service.GetService(log)
	app := api.InitCLIApp(ctx, log, service)

	err := app.Run(os.Args)
	if err != nil {
		return
	}
}
