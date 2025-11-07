package logger

import (
	"github.com/miles/rotki-demo/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var log *zap.Logger

// InitLogger 初始化日志记录器
func InitLogger(cfg *config.LogConfig) error {
	var zapConfig zap.Config

	if cfg.Output == "stdout" {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
		zapConfig.OutputPaths = []string{cfg.FilePath}
	}

	// 设置日志级别
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

// GetLogger 返回日志记录器实例
func GetLogger() *zap.Logger {
	if log == nil {
		// 如果未初始化，则回退到默认日志记录器
		log, _ = zap.NewDevelopment()
	}
	return log
}

// Info 记录信息消息
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Debug 记录调试消息
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Error 记录错误消息
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Warn 记录警告消息
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Fatal 记录致命消息并退出
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
	os.Exit(1)
}

// Sync 刷新任何缓冲的日志条目
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}
