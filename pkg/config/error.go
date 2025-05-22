package config

import (
	"GinBase/pkg/errs"
	"errors"
)

type ErrConfigNotFound interface {
	errs.WithCodeErr
	LackConfigName() string
}

type errConfigNotFound struct {
	errs.WithCodeErr
	lackConfigName string
}

func NewErrConfigNotFound(err error, lackConfigName string) ErrConfigNotFound {
	var wrappedErr errs.WithCodeErr
	if errors.As(err, &wrappedErr) && errs.IsErrorCode(wrappedErr, errs.ErrEnvironmentConfig) {
		return &errConfigNotFound{
			WithCodeErr:    wrappedErr,
			lackConfigName: lackConfigName,
		}

	}

	return &errConfigNotFound{
		WithCodeErr:    errs.WrapCodeError(errs.ErrEnvironmentConfig, err),
		lackConfigName: lackConfigName,
	}
}

func (e *errConfigNotFound) LackConfigName() string {
	return e.lackConfigName
}

type ErrConfigInvalid interface {
	errs.WithCodeErr
	InvalidConfig() string
}

type errConfigInvalid struct {
	errs.WithCodeErr
	invalidConfig string
}

func NewErrConfigInvalid(err error, invalidConfig string) ErrConfigInvalid {
	var wrappedErr errs.WithCodeErr
	if errors.As(err, &wrappedErr) && errs.IsErrorCode(err, errs.ErrEnvironmentConfig) {
		return &errConfigInvalid{
			WithCodeErr:   wrappedErr,
			invalidConfig: invalidConfig,
		}
	}

	return &errConfigInvalid{
		WithCodeErr:   errs.WrapCodeError(errs.ErrEnvironmentConfig, err),
		invalidConfig: invalidConfig,
	}
}

func (e *errConfigInvalid) InvalidConfig() string {
	return e.invalidConfig
}
