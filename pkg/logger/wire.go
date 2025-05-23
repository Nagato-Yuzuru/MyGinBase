//go:build wireinject
// +build wireinject

package logger

import (
	"github.com/google/wire"
	"terraqt.io/bedrock-go/pkg/config"
)

var LoggerSet = wire.NewSet(
	ProvideZapLogger,
)

func InitializeLogger() (Logger, error) {
	wire.Build(
		LoggerSet,
		config.ConfigSet,
	)
	return nil, nil
}
