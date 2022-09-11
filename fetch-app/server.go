package fetch_app

import (
	"fetch-app/internal/api"
	v1 "fetch-app/internal/api/v1"
	"fetch-app/internal/config"
	"fetch-app/internal/service"
	"fetch-app/internal/service/external"
	"github.com/gin-gonic/gin"
)

type Server struct {
	gin *gin.Engine
}

func NewServer(cfg *config.AppConfig) *Server {

	router := gin.Default()

	commodityService := service.NewCommodityService(cfg)

	currencyService := external.NewCurrencyService(cfg)

	commodityHandler := v1.NewCommodityController(cfg, commodityService, currencyService)

	registerHandlers(router, &api.HealthCheck{}, commodityHandler)

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
