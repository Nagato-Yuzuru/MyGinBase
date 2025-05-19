package config

import "sync"

var (
	loadOnce sync.Once
	Cfg      *Config
)

type ProvideParam struct {
	ConfigNames []string
	ConfigType  string
	ConfigPaths []string
	EnvPrefix   string
	Defaults    map[string]any
}

func ProvideConfig(
	l Loader,
	p ProvideParam,
) *Config {
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
	return Cfg
}
