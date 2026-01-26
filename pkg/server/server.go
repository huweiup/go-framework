package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release, test
}

type Server struct {
	Engine *gin.Engine
	Config Config
}

func New(cfg Config) *Server {
	if cfg.Mode != "" {
		gin.SetMode(cfg.Mode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	return &Server{
		Engine: r,
		Config: cfg,
	}
}

func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.Config.Port)
	return s.Engine.Run(addr)
}
