package config

import (
	"os"
	"strings"
	"sync"
)

var (
	loadOnce sync.Once
	Cfg      *Config
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

func GetConfigParam() configParam {
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

func ProvideConfig(
	l Loader,
	p configParam,
) Config {
	loadOnce.Do(
		func() {
			l.SetLoaderParams(p.ConfigNames, p.ConfigType, p.ConfigPaths, p.EnvPrefix, p.Defaults)

			var err error
			Cfg, err = l.Load()
			if err != nil {
				panic(err)
			}
		},
	)
	return *Cfg
}
