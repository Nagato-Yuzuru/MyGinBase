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
	ProvideConfig,
	wire.FieldsOf(new(Config), "Logger"),
)

func InitializeConfig() Config {
	wire.Build(ConfigSet)
	return Config{}
}
