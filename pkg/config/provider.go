package config

import (
	"os"
	"reflect"
	"strings"
	"sync"
)

var (
	loadOnce sync.Once
	Cfg      any
)

var (
	configsCacheMu  sync.Mutex
	configsCache    = make(map[reflect.Type]any)        // Cache for loaded configurations
	loadOncePerType = make(map[reflect.Type]*sync.Once) // Ensure single load per type
)

type configParam struct {
	ConfigNames []string
	ConfigType  string
	ConfigPaths []string
	EnvPrefix   string
	Defaults    map[string]any
}

const ENV = "ENV"

var (
	configParamOnce sync.Once
	param           configParam
)

func getConfigParam() configParam {
	configParamOnce.Do(
		func() {
			env, exist := os.LookupEnv(ENV)

			if !exist || env == "" {
				env = "PROD"
			}
			env = strings.ToLower(env)

			param = configParam{
				ConfigNames: []string{"base", env},
				ConfigType:  "yaml",
				ConfigPaths: []string{"./config"},
				EnvPrefix:   "",
				Defaults:    make(map[string]any),
			}
		},
	)
	return param
}

func ProvideGenericsConfig[T any](
	l Loader[T],
) T {
	var zeroT T
	configType := reflect.TypeOf(zeroT)

	configsCacheMu.Lock()

	once, ok := loadOncePerType[configType]
	if !ok {
		once = &sync.Once{}
		loadOncePerType[configType] = once
	}

	configsCacheMu.Unlock()

	once.Do(
		func() {
			param := getConfigParam()
			l.SetLoaderParams(
				param.ConfigNames,
				param.ConfigType,
				param.ConfigPaths,
				param.EnvPrefix,
				param.Defaults,
			)

			cfg, err := l.Load()
			if err != nil {
				panic(err)
			}

			if err := l.Valid(cfg); err != nil {
				panic(err)
			}

			configsCacheMu.Lock()
			configsCache[configType] = cfg
			configsCacheMu.Unlock()
		},
	)

	configsCacheMu.Lock()
	cfg, ok := configsCache[configType]
	configsCacheMu.Unlock()
	if !ok {
		panic("once done but config not found in cache")
	}

	cfgPtr, ok := cfg.(*T)
	if !ok {
		panic("config type mismatch")
	}

	return *cfgPtr
}

func ProvideConfig() Config {
	return ProvideGenericsConfig[Config](NewViperLoader[Config]())
}
