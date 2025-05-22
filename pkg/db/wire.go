//go:build wireinject
// +build wireinject

package db

import (
	"GinBase/pkg/config"
	"github.com/google/wire"
)

var PoolSet = wire.NewSet(
	ProvidePostgresPool,
)

func InitializePool() (Pool, error) {
	wire.Build(
		config.ConfigSet,
		PoolSet,
	)

	return nil, nil
}
