package database

import (
	"fmt"
	"time"

	"github.com/miles/rotki-demo/internal/config"
	"github.com/miles/rotki-demo/internal/logger"
	"github.com/miles/rotki-demo/internal/models"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DB 是全局数据库实例
var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase(cfg *config.DatabaseConfig) error {
	dsn := cfg.GetDSN()

	// 配置 GORM 日志记录器
	logLevel := gormlogger.Info
	gormLog := gormlogger.New(
		&gormLoggerAdapter{},
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// 打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLog,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 获取底层 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// 设置连接池设置
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	logger.Info("Database connected successfully",
		zap.String("host", cfg.Host),
		zap.String("database", cfg.Database),
	)

	return nil
}

// AutoMigrate 运行数据库迁移
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	err := DB.AutoMigrate(
		&models.Wallet{},
		&models.Address{},
		&models.Chain{},
		&models.Token{},
		&models.AssetSnapshot{},
		&models.SyncJob{},
		&models.RPCNode{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("Database migrations completed successfully")
	return nil
}

// GetDB 返回数据库实例
func GetDB() *gorm.DB {
	return DB
}

// gormLoggerAdapter 将 zap 日志记录器适配到 GORM 日志记录器接口
type gormLoggerAdapter struct{}

func (l *gormLoggerAdapter) Printf(format string, args ...interface{}) {
	logger.Debug(fmt.Sprintf(format, args...))
}
