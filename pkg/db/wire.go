//go:build wireinject
// +build wireinject

package db

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5"
	"terraqt.io/colas/bedrock-go/pkg/config"
	"terraqt.io/colas/bedrock-go/pkg/logger"
)

var PoolSet = wire.NewSet(
	provideAdaptPool,
	providePostgresPool,
	provideDriver,
	provideTransaction,
)

func InitializePGPool() (PGPool, error) {
	wire.Build(
		config.InitializeConfig,
		wire.FieldsOf(new(config.Config), "Database"),
		logger.InitializeLogger,
		PoolSet,
	)

	return nil, nil
}

func InitializeDriver() (*sql.Driver, error) {
	wire.Build(
		config.InitializeConfig,
		wire.FieldsOf(new(config.Config), "Database"),
		logger.InitializeLogger,
		PoolSet,
	)
	return nil, nil
}

func InitializeTx(tx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	wire.Build(
		config.InitializeConfig,
		wire.FieldsOf(new(config.Config), "Database"),
		logger.InitializeLogger,
		PoolSet,
	)
	return nil, nil
}
