package service

import (
	"go.uber.org/zap"
)

type UserService interface {
}

type Config interface {
	Logger() *zap.Logger
}
