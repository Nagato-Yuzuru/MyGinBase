//go:build wireinject
// +build wireinject

package logger

import (
	"github.com/google/wire"
	"terraqt.io/colas/bedrock-go/pkg/config"
)

var LoggerSet = wire.NewSet(
	provideZapLogger,
)

func InitializeLogger() (Logger, error) {
	wire.Build(
		LoggerSet,
		config.InitializeConfig,
		wire.FieldsOf(new(config.Config), "Logger"),
	)
	return nil, nil
}
