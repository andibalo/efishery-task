package v1

import (
	voerrors "auth-app/internal/autherrors"
	"auth-app/internal/config"
	"auth-app/internal/constants"
	"auth-app/internal/request"
	"auth-app/internal/response"
	"auth-app/internal/service"
	"auth-app/internal/util"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type UserController struct {
	cfg         config.Config
	userService service.UserService
}

func NewUserController(cfg config.Config, userService service.UserService) *UserController {
	return &UserController{
		cfg:         cfg,
		userService: userService,
	}
}

func (h *UserController) AddRoutes(e *echo.Echo) {
	r := e.Group(constants.V1BasePath + constants.UserBasePath)

	r.POST("/", h.createUser)
	r.POST("/login", h.login)
	r.POST("/details", h.getUserTokenDetails)
}

func (h *UserController) login(c echo.Context) error {

	loginUserReq := &request.LoginUserRequest{}

	if err := c.Bind(loginUserReq); err != nil {
		h.cfg.Logger().Error("login: error binding login request", zap.Error(err))

		return c.JSON(http.StatusInternalServerError, nil)
	}

	code, token, err := h.userService.GetJWTByPhoneAndPassword(loginUserReq.Phone, loginUserReq.Password)

	if err != nil {
		h.cfg.Logger().Error("login: error on user service", zap.Error(err))
		return h.failedUserResponse(c, code, err, "")
	}

	resp := response.NewResponse(code, token)

	resp.SetResponseMessage("Successfully logged in")

	return c.JSON(http.StatusOK, resp)
}

func (h *UserController) createUser(c echo.Context) error {

	registerUserReq := &request.CreateUserRequest{}

	if err := c.Bind(registerUserReq); err != nil {
		h.cfg.Logger().Error("createUser: error binding create request", zap.Error(err))

		return c.JSON(http.StatusInternalServerError, nil)
	}

	code, user, err := h.userService.CreateUser(registerUserReq)

	if err != nil {
		h.cfg.Logger().Error("createUser: error on user service", zap.Error(err))
		return h.failedUserResponse(c, code, err, "error creating user")
	}

	resp := response.NewResponse(code, user)

	resp.SetResponseMessage("Successfully created user")

	return c.JSON(http.StatusOK, resp)
}

func (h *UserController) getUserTokenDetails(c echo.Context) error {

	getUserTokenDetailsReq := &request.GetUserTokenDetailsRequest{}

	if err := c.Bind(getUserTokenDetailsReq); err != nil {
		h.cfg.Logger().Error("getUserTokenDetails: error binding request", zap.Error(err))

		return c.JSON(http.StatusInternalServerError, nil)
	}

	details, err := util.ParseToken(getUserTokenDetailsReq.Token)

	if err != nil {
		h.cfg.Logger().Error("getUserTokenDetails: error parsing token", zap.Error(err))
		return h.failedUserResponse(c, response.ServerError, err, "")
	}

	resp := response.NewResponse(response.Success, details)

	resp.SetResponseMessage("Successfully get token details")

	return c.JSON(http.StatusOK, resp)
}

func (h *UserController) failedUserResponse(c echo.Context, code response.Code, err error, errorMsg string) error {
	if code == "" {
		code = voerrors.MapErrorsToCode(err)
	}

	resp := response.Wrapper{
		ResponseCode: code,
		Status:       code.GetStatus(),
		Message:      code.GetMessage(),
	}

	if errorMsg != "" {
		resp.SetResponseMessage(errorMsg)
	}

	return c.JSON(voerrors.MapErrorsToStatusCode(err), resp)
}
