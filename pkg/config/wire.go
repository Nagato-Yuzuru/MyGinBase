//go:build wireinject
// +build wireinject

package config

import (
	"github.com/google/wire"
)

var ConfigSet = wire.NewSet(
	NewViperLoader,
	wire.Bind(new(Loader), new(*ViperLoader)),
	provideConfig,
	wire.FieldsOf(new(Config), "Logger", "Database"),
)

func InitializeConfig() Config {
	wire.Build(ConfigSet)
	return Config{}
}
