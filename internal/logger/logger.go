package logger

import (
	"github.com/miles/rotki-demo/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var log *zap.Logger

// InitLogger initializes the logger
func InitLogger(cfg *config.LogConfig) error {
	var zapConfig zap.Config

	if cfg.Output == "stdout" {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
		zapConfig.OutputPaths = []string{cfg.FilePath}
	}

	// Set log level
	level := zapcore.InfoLevel
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	var err error
	log, err = zapConfig.Build()
	if err != nil {
		return err
	}

	return nil
}

// GetLogger returns the logger instance
func GetLogger() *zap.Logger {
	if log == nil {
		// Fallback to a default logger if not initialized
		log, _ = zap.NewDevelopment()
	}
	return log
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
	os.Exit(1)
}

// Sync flushes any buffered log entries
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}
