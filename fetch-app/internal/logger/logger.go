package logger

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

var MainLogger *zap.Logger = nil

var once sync.Once

func NewMainLoggerSingleton() *zap.Logger {
	once.Do(func() {
		logger, err := zap.NewProduction()
		if err != nil {
			logger.Error("logger initialization failed", zap.Any("error", err))
			panic(fmt.Sprintf("logger initialization failed %v", err))
		}
		logger.Info("logger started")
		MainLogger = logger
	})

	return MainLogger
}

func GetMainLogger() *zap.Logger {
	return MainLogger
}
