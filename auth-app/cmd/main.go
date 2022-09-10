package main

import (
	"auth-app"
	"auth-app/internal/config"
	"fmt"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	cfg := config.InitConfig()

	server := auth_app.NewServer(cfg)

	err = server.Start(cfg.ServerAddress())

	if err != nil {
		cfg.Logger().Fatal("Port already used")
	}
}
