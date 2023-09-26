package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

// InitLog - инициализирует логгер
func InitLog() *logrus.Logger {
	log := logrus.New()
	// Включаем отслеживание вызывающего кода,
	// чтобы записывать информацию о местоположении вызовов функций.
	log.SetReportCaller(true)
	log.Out = os.Stdout
	log.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s():%d", frame.Function, frame.Line), filename
		},
		DisableColors:  true,
		DisableSorting: false,
	}
	return log
}
