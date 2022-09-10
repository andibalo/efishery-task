package api

import (
	"github.com/labstack/echo/v4"
)

type Handler interface {
	AddRoutes(router *echo.Echo)
}
