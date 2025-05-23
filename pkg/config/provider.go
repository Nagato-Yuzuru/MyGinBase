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

func provideConfig(
	l Loader,
) Config {
	loadOnce.Do(
		func() {
			p := getConfigParam()
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
