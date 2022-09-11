package config

import (
	logger2 "fetch-app/internal/logger"
	"github.com/bluele/gcache"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

const envVarEnvironment = "ENV"

func InitConfig() *AppConfig {
	logger := logger2.NewMainLoggerSingleton()

	gc := gcache.New(10).Build()

	return &AppConfig{
		logger:      logger,
		environment: os.Getenv(envVarEnvironment),
		cache:       gc,
	}
}

type AppConfig struct {
	logger      *zap.Logger
	environment string
	cache       gcache.Cache
}

type Config interface {
	Logger() *zap.Logger
	ServerAddress() string
	CurrencyServiceAPIKey() string
	GCache() gcache.Cache
}

func (a *AppConfig) Logger() *zap.Logger {
	return a.logger
}

func (a *AppConfig) ServerAddress() string {
	return viper.GetString("SERVER_PORT")
}

func (a *AppConfig) CurrencyServiceAPIKey() string {
	return viper.GetString("CURRENCY_SERVICE_API_KEY")
}

func (a *AppConfig) GCache() gcache.Cache {
	return a.cache
}
