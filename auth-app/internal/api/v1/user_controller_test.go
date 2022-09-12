package v1

import (
	"auth-app/internal/config"
	"auth-app/internal/model"
	"auth-app/internal/response"
	mocks2 "auth-app/internal/service/mocks"
	"auth-app/internal/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserController_AddRoutes(t *testing.T) {
	type fields struct {
		getConfig func() config.Config
	}
	type args struct {
		router *echo.Echo
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "OK",
			fields: fields{
				getConfig: func() config.Config {

					return config.InitConfig()
				},
			},
			args: args{
				router: echo.New(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &UserController{
				cfg: tt.fields.getConfig(),
			}
			assert.NotPanics(t, func() {
				d.AddRoutes(tt.args.router)
			})
		})
	}
}

func TestUserController_login(t *testing.T) {
	type fields struct {
		getConfig func() config.Config
	}

	validator := util.GetNewValidator()

	e := echo.New()

	tests := []struct {
		name                       string
		req                        string
		fields                     fields
		inGetJWTByPhoneAndPassword func(mockUserService *mocks2.UserService)
		responseCode               response.Code
	}{
		{
			name: "Error invalid request body",
			req:  getInvalidLoginJson(),
			fields: fields{
				getConfig: func() config.Config {

					return config.InitConfig()
				},
			},
			inGetJWTByPhoneAndPassword: func(mockUserService *mocks2.UserService) {
				mockUserService.On("GetJWTByPhoneAndPassword", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(response.Success, "jwt", nil)
			},
			responseCode: response.BadRequest,
		},
		{
			name: "Error fetching user",
			req:  getValidLoginJson(),
			fields: fields{
				getConfig: func() config.Config {

					return config.InitConfig()
				},
			},
			inGetJWTByPhoneAndPassword: func(mockUserService *mocks2.UserService) {
				mockUserService.On("GetJWTByPhoneAndPassword", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(response.ServerError, "", errors.New("some error"))
			},
			responseCode: response.ServerError,
		},
		{
			name: "OK",
			req:  getValidLoginJson(),
			fields: fields{
				getConfig: func() config.Config {

					return config.InitConfig()
				},
			},
			inGetJWTByPhoneAndPassword: func(mockUserService *mocks2.UserService) {
				mockUserService.On("GetJWTByPhoneAndPassword", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(response.Success, "jwt", nil)
			},
			responseCode: response.Success,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodPost, "/v1/user/login", strings.NewReader(tt.req))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			mockUserService := &mocks2.UserService{}

			tt.inGetJWTByPhoneAndPassword(mockUserService)

			h := &UserController{
				cfg:         tt.fields.getConfig(),
				userService: mockUserService,
				validator:   validator,
			}

			h.login(c)

			resBody, _ := io.ReadAll(rec.Body)

			var resBodyWrapper response.Wrapper
			_ = json.Unmarshal(resBody, &resBodyWrapper)
			resBodyCode := resBodyWrapper.ResponseCode
			assert.Equal(t, tt.responseCode, resBodyCode)
		})
	}
}

func TestUserController_createUser(t *testing.T) {
	type fields struct {
		getConfig func() config.Config
	}

	validator := util.GetNewValidator()

	e := echo.New()

	tests := []struct {
		name                       string
		req                        string
		fields                     fields
		inGetJWTByPhoneAndPassword func(mockUserService *mocks2.UserService)
		responseCode               response.Code
	}{
		{
			name: "Error invalid request body",
			req:  getInvalidLoginJson(),
			fields: fields{
				getConfig: func() config.Config {

					return config.InitConfig()
				},
			},
			inGetJWTByPhoneAndPassword: func(mockUserService *mocks2.UserService) {
				mockUserService.On("CreateUser", mock.Anything).Return(response.Success, model.User{}, nil)
			},
			responseCode: response.BadRequest,
		},
		{
			name: "Error creating user",
			req:  getValidCreateUserJson(),
			fields: fields{
				getConfig: func() config.Config {

					return config.InitConfig()
				},
			},
			inGetJWTByPhoneAndPassword: func(mockUserService *mocks2.UserService) {
				mockUserService.On("CreateUser", mock.Anything).Return(response.ServerError, model.User{}, errors.New("some error"))
			},
			responseCode: response.ServerError,
		},
		{
			name: "OK",
			req:  getValidCreateUserJson(),
			fields: fields{
				getConfig: func() config.Config {

					return config.InitConfig()
				},
			},
			inGetJWTByPhoneAndPassword: func(mockUserService *mocks2.UserService) {
				mockUserService.On("CreateUser", mock.Anything).Return(response.Success, model.User{}, nil)
			},
			responseCode: response.Success,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodPost, "/v1/user/", strings.NewReader(tt.req))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			mockUserService := &mocks2.UserService{}

			tt.inGetJWTByPhoneAndPassword(mockUserService)

			h := &UserController{
				cfg:         tt.fields.getConfig(),
				userService: mockUserService,
				validator:   validator,
			}

			h.createUser(c)

			resBody, _ := io.ReadAll(rec.Body)

			var resBodyWrapper response.Wrapper
			_ = json.Unmarshal(resBody, &resBodyWrapper)
			resBodyCode := resBodyWrapper.ResponseCode
			assert.Equal(t, tt.responseCode, resBodyCode)
		})
	}
}

func TestUserController_getUserTokenDetails(t *testing.T) {
	type fields struct {
		getConfig func() config.Config
	}

	validator := util.GetNewValidator()

	e := echo.New()

	tests := []struct {
		name         string
		req          string
		fields       fields
		responseCode response.Code
	}{
		{
			name: "Error invalid request body",
			req:  getInvalidLoginJson(),
			fields: fields{
				getConfig: func() config.Config {

					return config.InitConfig()
				},
			},
			responseCode: response.BadRequest,
		},
		{
			name: "OK",
			req:  getValidGetTokenDetailsJson(),
			fields: fields{
				getConfig: func() config.Config {

					return config.InitConfig()
				},
			},
			responseCode: response.Success,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodPost, "/v1/user/", strings.NewReader(tt.req))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			mockUserService := &mocks2.UserService{}

			h := &UserController{
				cfg:         tt.fields.getConfig(),
				userService: mockUserService,
				validator:   validator,
			}

			h.getUserTokenDetails(c)

			resBody, _ := io.ReadAll(rec.Body)

			var resBodyWrapper response.Wrapper
			_ = json.Unmarshal(resBody, &resBodyWrapper)
			resBodyCode := resBodyWrapper.ResponseCode
			assert.Equal(t, tt.responseCode, resBodyCode)
		})
	}
}

func getValidGetTokenDetailsJson() string {

	user := model.User{
		Name:       "andi",
		Phone:      "0000",
		Role:       "admin",
		Password:   "test",
		Timestampz: "123232",
	}
	token, _ := util.GenerateToken(user)

	return fmt.Sprintf(`{"token": "%s"}`, token)
}

func getValidCreateUserJson() string {
	return `{"name": "andi", "phone": "0000", "role": "admin"}`
}

func getValidLoginJson() string {
	return `{"phone": "0000", "password": "abcd"}`
}

func getInvalidLoginJson() string {
	return `{"name": "0000", "password": "abcd"}`
}
