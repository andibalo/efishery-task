package logger

import (
	"auth-app/internal/constants"
	"fmt"
	"sync"

	"github.com/labstack/echo/v4"

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

func ContextLogger(c echo.Context) *zap.Logger {
	logger := GetMainLogger()

	ctxRqID := c.Request().Header.Get(constants.HeaderRequestID)

	if ctxRqID != "" {
		return logger.With(zap.String(constants.HeaderRequestID, ctxRqID))
	}

	return logger
}
