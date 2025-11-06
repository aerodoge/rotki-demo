package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	DeBank   DeBankConfig   `mapstructure:"debank"`
	Sync     SyncConfig     `mapstructure:"sync"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	CacheTTL int    `mapstructure:"cache_ttl"`
}

type DeBankConfig struct {
	APIKey    string          `mapstructure:"api_key"`
	BaseURL   string          `mapstructure:"base_url"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	CacheTTL  int             `mapstructure:"cache_ttl"`
	Timeout   int             `mapstructure:"timeout"`
}

type RateLimitConfig struct {
	RequestsPerSecond int `mapstructure:"requests_per_second"`
	Burst             int `mapstructure:"burst"`
}

type SyncConfig struct {
	Enabled   bool `mapstructure:"enabled"`
	Interval  int  `mapstructure:"interval"`
	BatchSize int  `mapstructure:"batch_size"`
}

type LogConfig struct {
	Level    string `mapstructure:"level"`
	Output   string `mapstructure:"output"`
	FilePath string `mapstructure:"file_path"`
}

// LoadConfig loads configuration from file
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("redis.cache_ttl", 300)
	viper.SetDefault("debank.cache_ttl", 60)
	viper.SetDefault("debank.timeout", 30)
	viper.SetDefault("sync.enabled", true)
	viper.SetDefault("sync.interval", 300)
	viper.SetDefault("sync.batch_size", 10)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.output", "stdout")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.Charset,
	)
}

// GetRedisAddr returns the Redis address
func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetCacheTTL returns cache TTL as duration
func (c *RedisConfig) GetCacheTTL() time.Duration {
	return time.Duration(c.CacheTTL) * time.Second
}

// GetTimeout returns timeout as duration
func (c *DeBankConfig) GetTimeout() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}

// GetSyncInterval returns sync interval as duration
func (c *SyncConfig) GetSyncInterval() time.Duration {
	return time.Duration(c.Interval) * time.Second
}
