package fetch_app

import (
	"fetch-app/internal/api"
	"fetch-app/internal/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	gin *gin.Engine
}

func NewServer(cfg *config.AppConfig) *Server {

	router := gin.Default()

	registerHandlers(router, &api.HealthCheck{})

	return &Server{
		gin: router,
	}
}

func (s *Server) Start(addr string) error {
	return s.gin.Run(addr)
}

func registerHandlers(g *gin.Engine, handlers ...api.Handler) {
	for _, handler := range handlers {
		handler.AddRoutes(g)
	}
}
