package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Loader interface {
	Load() (*Config, error)
	Valid(cfg *Config) error
	SetLoaderParams(
		configNames []string,
		configType string,
		configPaths []string,
		envPrefix string,
		defaults map[string]any,
	)
}

type ViperLoader struct {
	*viper.Viper
	*validator.Validate

	configNames []string
	configType  string
	configPaths []string
	envPrefix   string
	defaults    map[string]any
}

func (l *ViperLoader) SetLoaderParams(
	configNames []string,
	configType string,
	configPaths []string,
	envPrefix string,
	defaults map[string]any,
) {

	l.configNames = configNames
	l.configType = configType
	l.configPaths = configPaths
	l.envPrefix = envPrefix
	l.defaults = defaults

	if envPrefix != "" {
		l.SetEnvPrefix(envPrefix)
	}

	godotenv.Load()

	l.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "__"))
	//l.AutomaticEnv()
	for k, v := range defaults {
		l.SetDefault(k, v)
	}
}

func NewViperLoader() *ViperLoader {
	return &ViperLoader{
		Viper: viper.New(), Validate: validator.New(),
	}
}

func (l *ViperLoader) bind(t reflect.Type, parent string) {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		// 跳过未导出字段 / 嵌入字段
		if f.PkgPath != "" {
			continue
		}

		tag := f.Tag.Get("mapstructure")
		if tag == "" || tag == "-" {
			continue
		}

		key := tag
		if parent != "" {
			key = parent + "." + tag
		}

		if f.Type.Kind() == reflect.Struct &&
			!(f.Type.PkgPath() == "time" && f.Type.Name() == "Duration") {
			l.bind(f.Type, key)
			continue
		}

		// 叶子：绑定环境变量
		_ = l.BindEnv(key)
	}
}

func (l *ViperLoader) Valid(cfg *Config) error {
	err := l.Struct(cfg)

	if err == nil {
		return nil
	}

	var invalidVE *validator.InvalidValidationError
	var validationErrs validator.ValidationErrors
	if errors.As(err, &invalidVE) {
		return invalidVE
	}

	if errors.As(err, &validationErrs) {
		var validationErrorMessages []string

		for _, fieldErr := range validationErrs {
			validationErrorMessages = append(
				validationErrorMessages,
				fmt.Sprintf(
					"fields '%s' valid faild: rule '%s', value is '%v'",
					fieldErr.Namespace(),
					fieldErr.Tag(),
					fieldErr.Value(),
				),
			)
		}
		return NewErrConfigInvalid(validationErrs, strings.Join(validationErrorMessages, "\n - "))
	}

	return err
}

func (l *ViperLoader) Load() (*Config, error) {

	var cfg Config

	l.bind(reflect.TypeOf(cfg), "")

	l.SetConfigType(l.configType)

	for _, path := range l.configPaths {
		l.AddConfigPath(path)
	}

	for i, name := range l.configNames {
		var err error
		l.SetConfigName(name)
		if i == 0 {
			err = l.ReadInConfig()
		} else {
			err = l.MergeInConfig()
		}

		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("read config error: %w", err)
		}
	}

	if err := l.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config unmarshal error: %w", err)
	}

	if err := l.Valid(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
