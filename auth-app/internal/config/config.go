package config

import (
	logger2 "auth-app/internal/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

const envVarEnvironment = "ENV"

func InitConfig() *AppConfig {
	logger := logger2.NewMainLoggerSingleton()

	return &AppConfig{
		logger:      logger,
		environment: os.Getenv(envVarEnvironment),
	}
}

type AppConfig struct {
	logger      *zap.Logger
	environment string
}

type Config interface {
	Logger() *zap.Logger
	ServerAddress() string
	UserDataFilePath() string
}

func (a *AppConfig) Logger() *zap.Logger {
	return a.logger
}

func (a *AppConfig) ServerAddress() string {
	return viper.GetString("SERVER_PORT")
}

func (a *AppConfig) UserDataFilePath() string {
	return viper.GetString("FILE_PATH")
}
