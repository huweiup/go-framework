package app

import (
	"fmt"

	"github.com/project/go-framework/pkg/config"
	"github.com/project/go-framework/pkg/database"
	"github.com/project/go-framework/pkg/logger"
	"github.com/project/go-framework/pkg/server"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Config holds the complete configuration for the application
type Config struct {
	App      AppConfig       `mapstructure:"app"`
	Log      logger.Config   `mapstructure:"log"`
	Database database.Config `mapstructure:"database"`
	Server   server.Config   `mapstructure:"server"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

// Application is the main container for the application components
type Application struct {
	Config Config
	Logger *zap.Logger
	DB     *gorm.DB
	Server *server.Server
}

// New creates a new Application instance
func New(configPath string) (*Application, error) {
	var cfg Config
	if err := config.Load(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize Logger
	log, err := logger.New(cfg.Log)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Initialize Database
	db, err := database.New(cfg.Database)
	if err != nil {
		// Log error but maybe don't fail if DB is optional?
		// For a framework, if DB config is present but fails, it should probably fail.
		// If DB config is empty, maybe skip?
		// For now, let's assume if it fails, it's fatal.
		log.Error("failed to initialize database", zap.Error(err))
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize Server
	srv := server.New(cfg.Server)

	// Add Logger Middleware
	srv.Engine.Use(server.LoggerMiddleware(log))

	return &Application{
		Config: cfg,
		Logger: log,
		DB:     db,
		Server: srv,
	}, nil
}

// Run starts the HTTP server
func (a *Application) Run() error {
	a.Logger.Info("Starting application",
		zap.String("name", a.Config.App.Name),
		zap.String("version", a.Config.App.Version),
		zap.Int("port", a.Config.Server.Port),
	)
	return a.Server.Run()
}
