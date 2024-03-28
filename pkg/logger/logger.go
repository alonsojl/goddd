package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

func New() (*logrus.Logger, error) {
	file, err := os.OpenFile(os.Getenv("LOG_PATH"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	formatter := &logrus.TextFormatter{
		TimestampFormat: "02-01-2006 15:04:05",
		FullTimestamp:   true,
		DisableColors:   false,
		CallerPrettyfier: func(frame *runtime.Frame) (string, string) {
			var (
				function = ""
				file     = fmt.Sprintf("%s:%d", path.Base(frame.File), frame.Line)
			)
			return function, file
		},
	}
	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		return nil, err
	}

	log := logrus.New()
	log.SetFormatter(formatter)
	log.SetReportCaller(true)
	log.SetLevel(level)
	log.SetOutput(file)

	return log, nil
}
