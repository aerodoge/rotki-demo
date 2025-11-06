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

// DB is the global database instance
var DB *gorm.DB

// InitDatabase initializes the database connection
func InitDatabase(cfg *config.DatabaseConfig) error {
	dsn := cfg.GetDSN()

	// Configure GORM logger
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

	// Open database connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLog,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
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

// AutoMigrate runs database migrations
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

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// gormLoggerAdapter adapts zap logger to GORM logger interface
type gormLoggerAdapter struct{}

func (l *gormLoggerAdapter) Printf(format string, args ...interface{}) {
	logger.Debug(fmt.Sprintf(format, args...))
}
