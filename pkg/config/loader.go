package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Loader interface {
	Load() (*Config, error)
	Valid(cfg *Config) error
}

type ViperLoader struct {
	*viper.Viper
	*validator.Validate

	configName  string
	configType  string
	configPaths []string
	envPrefix   string
	defaults    map[string]interface{}
}

func NewViperLoader(configName string, configType string, configPaths []string, envPrefix string, defaults map[string]interface{}) *ViperLoader {
	return &ViperLoader{
		Viper: viper.New(), Validate: validator.New(),
		configName: configName, configType: configType, configPaths: configPaths, envPrefix: envPrefix, defaults: defaults}
}

func (l *ViperLoader) bind(t reflect.Type, parent string) {
	// 若 t 是指针，取 Elem
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

		// 拼接完整 key：parent.child
		key := tag
		if parent != "" {
			key = parent + "." + tag
		}

		// 继续向下递归（排除 time.Duration 等特殊结构体）
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
	err := l.Struct(&cfg)

	if err == nil {
		return nil
	}

	var invalidVE *validator.InvalidValidationError
	var validationErrs validator.ValidationErrors

	if errors.As(err, &invalidVE) {
		return fmt.Errorf("invalid params %w", err)
	}

	if errors.As(err, &validationErrs) {
		var validationErrorMessages []string

		for _, fieldErr := range validationErrs {
			validationErrorMessages = append(validationErrorMessages,
				fmt.Sprintf("fields '%s' valid faild: rule '%s', value is '%v'", fieldErr.Namespace(), fieldErr.Tag(), fieldErr.Value()))
		}
		return fmt.Errorf("invalid params:\n - %s", strings.Join(validationErrorMessages, "\n - "))
	}

	return err
}

func (l *ViperLoader) init() {
	l.Viper.SetConfigName(l.configName)
	l.Viper.SetConfigType(l.configType)

	for _, path := range l.configPaths {
		l.Viper.AddConfigPath(path)
	}

}

func (l *ViperLoader) Load() (*Config, error) {
	l.init()

	if err := l.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("config load error: %w", err)
		}
	}

	if l.envPrefix != "" {
		l.SetEnvPrefix(l.envPrefix)
	}

	l.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "__"))
	l.AutomaticEnv()
	for k, v := range l.defaults {
		l.SetDefault(k, v)
	}

	var cfg Config

	l.bind(reflect.TypeOf(cfg), "")

	if err := l.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config unmarshal error: %w", err)
	}

	if err := l.Valid(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
