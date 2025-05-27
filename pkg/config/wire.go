//go:build wireinject
// +build wireinject

package config

import (
	"github.com/google/wire"
)

func InitializeConfig() Config {
	wire.Build(
		ProvideConfig,
	)
	return Config{}
}
