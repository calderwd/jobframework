package logger

import (
	"log"

	"go.uber.org/zap"

	"sync"
)

var once sync.Once
var instance loggerStdout

type Logger interface {
	Error(err error, msg ...string)
	Warn(msg ...string)
	Info(msg ...string)
	Debug(msg ...string)
}

type loggerStdout struct {
	logger *zap.Logger
}

func (l loggerStdout) Error(err error, msg ...string) {
	l.logger.Sugar().Error(err, msg)
}

func (l loggerStdout) Warn(msg ...string) {
	l.logger.Sugar().Warn(msg)
}

func (l loggerStdout) Info(msg ...string) {
	l.logger.Sugar().Info(msg)
}

func (l loggerStdout) Debug(msg ...string) {
	l.logger.Sugar().Debug(msg)
}

func GetInstance() Logger {
	once.Do(func() {
		logger, err := zap.NewProduction()
		if err != nil {
			log.Fatal("failed to create logger")
		}

		instance = loggerStdout{
			logger: logger,
		}
	})
	return instance
}
