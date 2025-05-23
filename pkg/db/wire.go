//go:build wireinject
// +build wireinject

package db

import (
	"github.com/google/wire"
	"terraqt.io/bedrock-go/pkg/config"
	"terraqt.io/bedrock-go/pkg/logger"
)

var PoolSet = wire.NewSet(
	ProvidePostgresPool,
)

func InitializePGPool() (PGPool, error) {
	wire.Build(
		logger.LoggerSet,
		config.ConfigSet,
		PoolSet,
	)

	return nil, nil
}
