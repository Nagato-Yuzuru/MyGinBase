package db

import "time"

// DatabaseConfig 数据库连接配置
type DatabaseConfig struct {
	Name            string        `mapstructure:"name" validate:"required"`
	Driver          string        `mapstructure:"driver" validate:"required,oneof=postgres mysql sqlite"`
	Host            string        `mapstructure:"host" validate:"required"`
	Port            int32         `mapstructure:"port" validate:"required"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database" validate:"required"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int32         `mapstructure:"max_open_conns"`
	MaxIdleConns    int32         `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	DebugSQL        bool          `mapstructure:"debug_sql"`
}

type ConfigGetter interface {
	GetDbConfig() DatabaseConfig
}
