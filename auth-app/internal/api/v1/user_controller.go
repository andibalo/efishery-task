package v1

import (
	voerrors "auth-app/internal/autherrors"
	"auth-app/internal/config"
	"auth-app/internal/constants"
	"auth-app/internal/request"
	"auth-app/internal/response"
	"auth-app/internal/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type UserController struct {
	cfg         config.Config
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (h *UserController) AddRoutes(e *echo.Echo) {
	r := e.Group(constants.V1BasePath + constants.UserBasePath)

	r.POST("/", h.createUser)
}

func (h *UserController) createUser(c echo.Context) error {

	registerUserReq := &request.CreateUserRequest{}

	if err := c.Bind(registerUserReq); err != nil {
		h.cfg.Logger().Error("createUser: error binding create request", zap.Error(err))

		return c.JSON(http.StatusInternalServerError, nil)
	}

	code, err := h.userService.CreateUser(registerUserReq)

	if err != nil {
		h.cfg.Logger().Error("createUser: error on user service", zap.Error(err))
		return h.failedUserResponse(c, code, err, "")
	}

	resp := response.NewResponse(code, nil)

	resp.SetResponseMessage("Successfully created user")

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
