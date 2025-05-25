//go:build wireinject
// +build wireinject

package config

import (
	"github.com/google/wire"
)

var ConfigSet = wire.NewSet(
	newViperLoader,
	wire.Bind(new(Loader), new(*ViperLoader)),
	provideConfig,
)

func InitializeConfig() Config {
	wire.Build(ConfigSet)
	return Config{}
}
