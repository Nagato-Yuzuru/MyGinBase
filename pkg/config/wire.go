//go:build wireinject
// +build wireinject

package config

import (
	"github.com/google/wire"
)

var ConfigSet = wire.NewSet(
	NewViperLoader,
	GetConfigParam,
	wire.Bind(new(Loader), new(*ViperLoader)),
	config.ProvideConfig,
	wire.FieldsOf(new(config.Config), "Logger", "DatabaseD"),
)

func InitializeConfig() *Config {
	wire.Build(ConfigSet)
	return nil
}
