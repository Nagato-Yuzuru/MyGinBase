package config

import (
	"time"
)

// Config 是应用程序的完整配置结构
type Config struct {
	//App AppConfig `mapstructure:"app" validate:"required"`
	//Server    ServerConfig    `mapstructure:"server" validate:"required"`
	//Database  DatabaseConfig  `mapstructure:"database" validate:"required"`
	Logger LoggerConfig `mapstructure:"logger" validate:"required"`
	//Cache     CacheConfig     `mapstructure:"cache"`
	//Telemetry TelemetryConfig `mapstructure:"telemetry"`
}

// AppConfig 包含应用基础配置
type AppConfig struct {
	Name        string `mapstructure:"name" validate:"required"`
	Environment string `mapstructure:"environment" validate:"required,oneof=development testing staging production"`
	Version     string `mapstructure:"version"`
	Debug       bool   `mapstructure:"debug"`
}

// ServerConfig HTTP服务器配置
type ServerConfig struct {
	Host            string        `mapstructure:"host" validate:"required"`
	Port            int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	TrustedProxies  []string      `mapstructure:"trusted_proxies"`
}

// DatabaseConfig 数据库连接配置
type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver" validate:"required,oneof=postgres mysql sqlite"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database" validate:"required"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level  string `mapstructure:"level" validate:"required,oneof=debug info warn error dpanic panic fatal"`
	Format string `mapstructure:"format" validate:"required,oneof=json console"`
	//OutputPath string `mapstructure:"output_path"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Enabled    bool          `mapstructure:"enabled"`
	Type       string        `mapstructure:"type" validate:"omitempty,oneof=memory redis"`
	Addr       string        `mapstructure:"addr"`
	Password   string        `mapstructure:"password"`
	DB         int           `mapstructure:"db"`
	DefaultTTL time.Duration `mapstructure:"default_ttl"`
}

// TelemetryConfig 遥测配置
type TelemetryConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	OtelEndpoint string `mapstructure:"otel_endpoint"`
	ServiceName  string `mapstructure:"service_name"`
}

//// DefaultConfig 返回带有默认值的配置
//func DefaultConfig() *Config {
//	return &Config{
//		App: AppConfig{
//			Name:        "web-service",
//			Environment: "development",
//			Version:     "0.1.0",
//			Debug:       false,
//		},
//	Server: ServerConfig{
//		Host:            "0.0.0.0",
//		Port:            8080,
//		ReadTimeout:     5 * time.Second,
//		WriteTimeout:    10 * time.Second,
//		ShutdownTimeout: 30 * time.Second,
//		TrustedProxies:  []string{"127.0.0.1", "::1"},
//	},
//	Database: DatabaseConfig{
//		Driver:          "postgres",
//		Host:            "localhost",
//		Port:            5432,
//		Username:        "postgres",
//		Password:        "",
//		Database:        "app_db",
//		SSLMode:         "disable",
//		MaxOpenConns:    10,
//		MaxIdleConns:    5,
//		ConnMaxLifetime: 5 * time.Minute,
//	},
//	Logger: LoggerConfig{
//		Level:      "info",
//		Format:     "json",
//		OutputPath: "stdout",
//	},
//	Cache: CacheConfig{
//		Enabled:    false,
//		Type:       "memory",
//		DefaultTTL: 5 * time.Minute,
//	},
//	Telemetry: TelemetryConfig{
//		Enabled:      true,
//		OtelEndpoint: "localhost:4317",
//		ServiceName:  "web-service",
//	},
//}
//}
