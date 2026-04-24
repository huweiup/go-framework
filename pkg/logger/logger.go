package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *zap.Logger

type Config struct {
	Mode       string `mapstructure:"mode"` // dev, debug, release
	Level      string `mapstructure:"level"`
	MaxSize    int    `mapstructure:"max_size"`    // in MB
	MaxBackups int    `mapstructure:"max_backups"` // number of files
	MaxAge     int    `mapstructure:"max_age"`     // in days
	Compress   bool   `mapstructure:"compress"`
	Encoding   string `mapstructure:"encoding"` // "json" or "console"
}

// New creates a new zap logger with the given configuration
func New(cfg Config) error {
	cfgLevel, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return err
	}

	if cfg.Mode == "" {
		cfg.Mode = "release"
	}

	if cfg.MaxSize == 0 {
		cfg.MaxSize = 100
	}
	if cfg.MaxBackups == 0 {
		cfg.MaxBackups = 10
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = 30
	}

	var encoder zapcore.Encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if cfg.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	infoWriter := &lumberjack.Logger{
		Filename:   "./logs/info.log",
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	errorWriter := &lumberjack.Logger{
		Filename:   "./logs/error.log",
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= cfgLevel && level < zapcore.ErrorLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.ErrorLevel
	})

	infoCore := zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel)
	errorCore := zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel)
	var core zapcore.Core
	if cfg.Mode == "dev" || cfg.Mode == "debug" {
		stdCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.InfoLevel)
		core = zapcore.NewTee(stdCore, infoCore, errorCore)
	} else {
		core = zapcore.NewTee(infoCore, errorCore)
	}

	log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return nil
}

func Log() *zap.Logger {
	return log
}

func Logsugar() *zap.SugaredLogger {
	return log.Sugar()
}
