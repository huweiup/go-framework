package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/huweiup/go-framework/pkg/logger"
	"github.com/opslead/gormzap"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type Source struct {
	Master string
	Slave  []string
}

type Config struct {
	Driver       string `mapstructure:"driver"` // mysql
	Source       Source `mapstructure:"source"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"` // 打开到数据库的最大连接数

}

var gormDB *gorm.DB
var once sync.Once

// New creates a new Gorm DB connection
func New(cfg Config) (e error) {
	once.Do(func() {

		var dialector gorm.Dialector

		switch cfg.Driver {
		case "mysql":
			dialector = mysql.Open(cfg.Source.Master)
		default:
			e = fmt.Errorf("unsupported db driver: %s", cfg.Driver)
			return
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
			e = fmt.Errorf("failed to connect db: %w", err)
			return
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
			e = fmt.Errorf("failed to register db resolver: %w", err)
			return
		}

		sqlDB, _ := db.DB()
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

		gormDB = db
	})
	return
}

func GetDB() *gorm.DB {
	return gormDB
}

func Close() {
	sqlDB, _ := gormDB.DB()
	err := sqlDB.Close()
	if err != nil {
		logger.Log().Error("failed to close db connection", zap.Error(err))
	}

}
