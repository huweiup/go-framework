package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	cfg := Config{
		Port: 8080,
		Mode: "test",
	}

	srv := New(cfg)
	assert.NotNil(t, srv)
	assert.NotNil(t, srv.Engine)
	assert.Equal(t, cfg, srv.Config)
	assert.Equal(t, gin.TestMode, gin.Mode())
}

func TestLoggerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger, _ := zap.NewDevelopment()

	r := gin.New()
	r.Use(LoggerMiddleware(logger))

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
