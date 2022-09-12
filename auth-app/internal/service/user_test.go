package service

import (
	"auth-app/internal/config"
	"auth-app/internal/model"
	"auth-app/internal/request"
	"auth-app/internal/response"
	"auth-app/internal/storage/mocks"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name               string
		req                *request.CreateUserRequest
		inFetchUserByPhone func(mockStorage *mocks.Storage)
		inCreateUser       func(mockStorage *mocks.Storage)
		responseCode       response.Code
		wantErr            bool
	}{
		{
			name: "Error fetcing user by phone",
			req:  getValidCreateUserReq(),
			inFetchUserByPhone: func(mockStorage *mocks.Storage) {
				mockStorage.On("FetchUserByPhone", mock.AnythingOfType("string")).Return(model.User{}, errors.New("some error"))
			},
			inCreateUser: func(mockStorage *mocks.Storage) {
				mockStorage.On("CreateUser", mock.AnythingOfType("model.User")).Return(nil)
			},
			responseCode: response.ServerError,
			wantErr:      true,
		},
		{
			name: "Error duplicate user",
			req:  getValidCreateUserReq(),
			inFetchUserByPhone: func(mockStorage *mocks.Storage) {
				mockStorage.On("FetchUserByPhone", mock.AnythingOfType("string")).Return(model.User{Name: "test"}, nil)
			},
			inCreateUser: func(mockStorage *mocks.Storage) {
				mockStorage.On("CreateUser", mock.AnythingOfType("model.User")).Return(nil)
			},
			responseCode: response.BadRequest,
			wantErr:      true,
		},
		{
			name: "Error creating user",
			req:  getValidCreateUserReq(),
			inFetchUserByPhone: func(mockStorage *mocks.Storage) {
				mockStorage.On("FetchUserByPhone", mock.AnythingOfType("string")).Return(model.User{}, nil)
			},
			inCreateUser: func(mockStorage *mocks.Storage) {
				mockStorage.On("CreateUser", mock.AnythingOfType("model.User")).Return(errors.New("some error"))
			},
			responseCode: response.ServerError,
			wantErr:      true,
		},
		{
			name: "Ok",
			req:  getValidCreateUserReq(),
			inFetchUserByPhone: func(mockStorage *mocks.Storage) {
				mockStorage.On("FetchUserByPhone", mock.AnythingOfType("string")).Return(model.User{}, nil)
			},
			inCreateUser: func(mockStorage *mocks.Storage) {
				mockStorage.On("CreateUser", mock.AnythingOfType("model.User")).Return(nil)
			},
			responseCode: response.Success,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appCfg := config.InitConfig()

			mockStorage := &mocks.Storage{}
			tt.inFetchUserByPhone(mockStorage)
			tt.inCreateUser(mockStorage)

			service := NewUserService(appCfg, mockStorage)

			code, _, err := service.CreateUser(tt.req)

			if tt.wantErr {
				assert.NotNil(t, err)

				return
			}

			assert.Equal(t, code, tt.responseCode)
		})
	}
}

func TestUserService_GetJWTByPhoneAndPassword(t *testing.T) {
	tests := []struct {
		name               string
		inFetchUserByPhone func(mockStorage *mocks.Storage)
		responseCode       response.Code
		wantErr            bool
	}{
		{
			name: "Error fetcing user by phone and password",
			inFetchUserByPhone: func(mockStorage *mocks.Storage) {
				mockStorage.On("FetchUserByPhoneAndPassword", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(model.User{}, errors.New("some error"))
			},
			responseCode: response.ServerError,
			wantErr:      true,
		},

		{
			name: "Ok",
			inFetchUserByPhone: func(mockStorage *mocks.Storage) {
				mockStorage.On("FetchUserByPhoneAndPassword", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(model.User{}, nil)
			},
			responseCode: response.Success,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appCfg := config.InitConfig()

			mockStorage := &mocks.Storage{}
			tt.inFetchUserByPhone(mockStorage)

			service := NewUserService(appCfg, mockStorage)

			code, _, err := service.GetJWTByPhoneAndPassword("test", "test")

			if tt.wantErr {
				assert.NotNil(t, err)

				return
			}

			assert.Equal(t, code, tt.responseCode)
		})
	}
}

func getValidCreateUserReq() *request.CreateUserRequest {

	return &request.CreateUserRequest{
		Name:  "Andi",
		Phone: "0932832",
		Role:  "admin",
	}
}
