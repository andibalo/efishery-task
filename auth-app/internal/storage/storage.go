package storage

import (
	"auth-app/internal/config"
	"auth-app/internal/model"
	"auth-app/internal/storage/repositories"
	"go.uber.org/zap"
)

type Store struct {
	logger         *zap.Logger
	userRepository UserRepo
}

type Storage interface {
	CreateUser(user model.User) error
	FetchAllUsers() (users []model.User, err error)
	FetchUserByPhoneAndPassword(phone, password string) (user model.User, err error)
}

func New(cfg config.Config) *Store {

	userRepo := repositories.NewUserRepository(cfg)

	store := &Store{
		logger:         cfg.Logger(),
		userRepository: userRepo,
	}

	return store
}

type UserRepo interface {
	SaveUser(user model.User) (err error)
	GetAllUsers() (users []model.User, err error)
	GetUserByPhoneAndPassword(phone, password string) (user model.User, err error)
}
