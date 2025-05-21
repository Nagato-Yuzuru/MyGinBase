//go:build wireinject
// +build wireinject

package inject

import (
	"GinBase/pkg/config"
	"github.com/google/wire"
	"os"
	"strings"
	"sync"
)

// var LoggerSet = wire.NewSet(
//
//	logger.ProvideZapLogger,
//	wire.Bind(new(logger.Logger), new(*logger.ZapLogger)),
//
// )

// func InitLogger() logger.Logger {
//
//		wire.Build(LoggerSet)
//		return nil
//	}

var (
	viperLoader     config.ViperLoader
	configParamOnce sync.Once
	ConfigParam     config.configParam
)

const ENV = "ENV"

func NewProvideParam() config.configParam {
	configParamOnce.Do(
		func() {
			env, exist := os.LookupEnv(ENV)

			if !exist || env == "" {
				env = "PROD"
			}
			env = strings.ToLower(env)

			ConfigParam = config.configParam{
				ConfigNames: []string{"base", env},
				ConfigType:  "yaml",
				ConfigPaths: []string{"./config"},
				EnvPrefix:   "",
				Defaults:    make(map[string]any),
			}
		},
	)
	return ConfigParam
}

var ConfigSet = wire.NewSet(
	config.NewViperLoader,
	NewProvideParam,
	wire.Bind(new(config.Loader), new(*config.ViperLoader)),
	config.ProvideConfig,
	wire.FieldsOf(new(config.Config), "Logger", "DatabaseD"),
)

func InitializeConfig() *config.Config {
	wire.Build(ConfigSet)
	return nil
}
