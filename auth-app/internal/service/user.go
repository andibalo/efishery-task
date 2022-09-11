package service

import (
	voerrors "auth-app/internal/autherrors"
	"auth-app/internal/model"
	"auth-app/internal/request"
	"auth-app/internal/response"
	"auth-app/internal/storage"
	"auth-app/internal/util"
	"errors"
	"github.com/sethvargo/go-password/password"
	"go.uber.org/zap"
	"time"
)

type userService struct {
	config Config
	store  storage.Storage
}

func NewUserService(config Config, store storage.Storage) *userService {

	return &userService{
		config: config,
		store:  store,
	}
}

func (s *userService) CreateUser(createUserReq *request.CreateUserRequest) (response.Code, model.User, error) {
	s.config.Logger().Info("CreateUser: creating user")

	var user model.User

	existingUser, err := s.store.FetchUserByPhone(createUserReq.Phone)

	if err != nil {
		if !errors.Is(err, voerrors.ErrNotFound) {
			s.config.Logger().Error("CreateUser: error fetching user by phone", zap.Error(err))
			return response.ServerError, user, err
		}
	}

	userIsExist := existingUser != model.User{}

	if userIsExist {
		s.config.Logger().Error("CreateUser: duplicate user", zap.Error(err))
		return response.BadRequest, user, err
	}

	pass, err := s.generatePassword()

	if err != nil {
		s.config.Logger().Error("CreateUser: error generating password", zap.Error(err))
		return response.ServerError, user, err
	}

	timestampz := time.Now()

	user = model.User{
		Name:       createUserReq.Name,
		Phone:      createUserReq.Phone,
		Role:       createUserReq.Role,
		Password:   pass,
		Timestampz: timestampz.Format("02 Jan 06 15:04 MST"),
	}

	err = s.store.CreateUser(user)

	if err != nil {
		s.config.Logger().Error("CreateUser: error creating user", zap.Error(err))
		return response.ServerError, model.User{}, err
	}

	return response.Success, user, nil
}

func (s *userService) GetJWTByPhoneAndPassword(phone, password string) (response.Code, string, error) {

	s.config.Logger().Info("GetJWTByPhoneAndPassword: getting jwt by phone and password")
	user, err := s.store.FetchUserByPhoneAndPassword(phone, password)

	if err != nil {
		if !errors.Is(err, voerrors.ErrNotFound) {
			s.config.Logger().Error("GetJWTByPhoneAndPassword: error getting jwt by phone and password", zap.Error(err))
			return response.ServerError, "", err
		}
	}

	token, err := util.GenerateToken(user)

	if err != nil {
		s.config.Logger().Error("GetJWTByPhoneAndPassword: error generating token", zap.Error(err))
		return response.ServerError, "", err
	}

	return response.Success, token, nil
}

func (s *userService) generatePassword() (pass string, err error) {
	pass, err = password.Generate(4, 1, 1, false, false)
	if err != nil {
		return "", err
	}

	return pass, nil
}
