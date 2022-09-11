package service

import (
	"auth-app/internal/model"
	"auth-app/internal/request"
	"auth-app/internal/response"

	"go.uber.org/zap"
)

type UserService interface {
	GetJWTByPhoneAndPassword(phone, password string) (response.Code, string, error)
	CreateUser(createUserReq *request.CreateUserRequest) (response.Code, model.User, error)
}

type Config interface {
	Logger() *zap.Logger
}
