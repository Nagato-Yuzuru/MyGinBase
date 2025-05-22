package config

import (
	"GinBase/pkg/errs"
	"errors"
)

type ErrConfigNotFound interface {
	errs.CodedError
	LackConfigName() string
}

type errConfigNotFound struct {
	errs.CodedError
	lackConfigName string
}

func NewErrConfigNotFound(err error, lackConfigName string) ErrConfigNotFound {
	var wrappedErr errs.CodedError
	if errors.As(err, &wrappedErr) && errs.IsErrorCode(wrappedErr, errs.ErrEnvironmentConfig) {
		return &errConfigNotFound{
			CodedError:     wrappedErr,
			lackConfigName: lackConfigName,
		}

	}

	return &errConfigNotFound{
		CodedError:     errs.WrapCodeError(errs.ErrEnvironmentConfig, err),
		lackConfigName: lackConfigName,
	}
}

func (e *errConfigNotFound) LackConfigName() string {
	return e.lackConfigName
}

type ErrConfigInvalid interface {
	errs.CodedError
	InvalidConfig() string
}

type errConfigInvalid struct {
	errs.CodedError
	invalidConfig string
}

func NewErrConfigInvalid(err error, invalidConfig string) ErrConfigInvalid {
	var wrappedErr errs.CodedError
	if errors.As(err, &wrappedErr) && errs.IsErrorCode(err, errs.ErrEnvironmentConfig) {
		return &errConfigInvalid{
			CodedError:    wrappedErr,
			invalidConfig: invalidConfig,
		}
	}

	return &errConfigInvalid{
		CodedError:    errs.WrapCodeError(errs.ErrEnvironmentConfig, err),
		invalidConfig: invalidConfig,
	}
}

func (e *errConfigInvalid) InvalidConfig() string {
	return e.invalidConfig
}
