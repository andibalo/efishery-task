package auth_app

import (
	"auth-app/internal/api"
	v1 "auth-app/internal/api/v1"
	"auth-app/internal/config"
	"auth-app/internal/service"
	"auth-app/internal/storage"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo *echo.Echo
}

func NewServer(cfg *config.AppConfig) *Server {

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid.New().String()
		},
	}))

	store := storage.New(cfg)

	userService := service.NewUserService(cfg, store)
	userHandler := v1.NewUserController(userService)

	registerHandlers(e, &api.HealthCheck{}, userHandler)

	return &Server{
		echo: e,
	}
}

func (s *Server) Start(addr string) error {
	return s.echo.Start(":" + addr)
}

func registerHandlers(e *echo.Echo, handlers ...api.Handler) {
	for _, handler := range handlers {
		handler.AddRoutes(e)
	}
}
