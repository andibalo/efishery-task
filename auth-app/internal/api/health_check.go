package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	healthCheckPath = "/health"
)

// HealthCheck is a standard, simple health check
type HealthCheck struct{}

// AddRoutes adds the routers for this API to the provided router (or subrouter)
func (h *HealthCheck) AddRoutes(e *echo.Echo) {
	e.GET(healthCheckPath, h.handler)
}

func (h *HealthCheck) handler(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
