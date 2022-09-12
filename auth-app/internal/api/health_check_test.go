package api

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHealthCheck(t *testing.T) {
	e := echo.New()
	healthCheck := HealthCheck{}
	healthCheck.AddRoutes(e)

	// function call
	resp := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/health", http.NoBody)
	assert.NoError(t, err)

	e.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code, "status codes")

	payload, _ := io.ReadAll(resp.Result().Body)
	assert.True(t, len(payload) > 0, "payload")
}
