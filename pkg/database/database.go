package database

import (
	"fmt"
	"time"

	"github.com/huweiup/go-framework/pkg/logger"
	"github.com/opslead/gormzap"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type Source struct {
	Master string
	Slave  []string
}

type Config struct {
	Driver string `mapstructure:"driver"` // mysql, sqlite
	Source Source `mapstructure:"source"`
}

var gormDB *gorm.DB

// New creates a new Gorm DB connection
func New(cfg Config) error {
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.Source.Master)
	case "sqlite":
		dialector = sqlite.Open(cfg.Source.Master)
	default:
		return fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	gormZapLogger := gormzap.NewGormZapLogger(gormzap.GormZapLoggerConfig{
		Logger:                    logger.Log(),
		LogLevel:                  gormLogger.Info,        // 日志级别
		SlowThreshold:             200 * time.Millisecond, // 慢查询阈值
		IgnoreRecordNotFoundError: true,                   // 忽略记录未找到错误
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormZapLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}
	sources := []gorm.Dialector{mysql.Open(cfg.Source.Master)}
	replices := make([]gorm.Dialector, 0, len(cfg.Source.Slave))
	for _, slave := range cfg.Source.Slave {
		replices = append(replices, mysql.Open(slave))
	}

	if err := db.Use(dbresolver.Register(dbresolver.Config{
		Sources:           sources,
		Replicas:          replices,
		Policy:            dbresolver.RandomPolicy{},
		TraceResolverMode: true,
	})); err != nil {
		logger.Log().Error("failed to register db resolver", zap.Error(err))
		return fmt.Errorf("failed to register db resolver: %w", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)

	gormDB = db
	return nil
}

func GetDB() *gorm.DB {
	return gormDB
}

func Close() {
	sqlDB, _ := gormDB.DB()
	err := sqlDB.Close()
	if err != nil {
		logger.Log().Error("failed to close database connection", zap.Error(err))
	}

}
