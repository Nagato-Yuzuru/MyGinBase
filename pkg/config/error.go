package config

import (
	"GinBase/pkg/errs"
	"fmt"
)

type ErrConfigNotFound interface {
	errs.WithCodeErr
	LackConfigName() string
	Unwrap() error
}

type errConfigNotFound struct {
	lackConfigName string
	wrappedErr     error
}

func NewErrConfigNotFound(lackConfigName string, wrappedErr error) ErrConfigNotFound {
	return &errConfigNotFound{lackConfigName: lackConfigName, wrappedErr: wrappedErr}
}

func (e *errConfigNotFound) Unwrap() error {
	return e.wrappedErr
}

func (e *errConfigNotFound) Error() string {
	return fmt.Sprintf("config not found: %s", e.lackConfigName)
}

func (e *errConfigNotFound) Code() errs.Code {
	return errs.ErrNotFound
}

func (e *errConfigNotFound) LackConfigName() string {
	return e.lackConfigName
}

type ErrConfigInvalid interface {
	errs.WithCodeErr
	InvalidConfig() string
	Unwrap() error
}

type errConfigInvalid struct {
	invalidConfig string
	wrappedErr    error
}

func (e *errConfigInvalid) Code() errs.Code {
	return errs.ErrInvalidParam
}

func NewErrConfigInvalid(invalidConfig string, wrappedErr error) ErrConfigInvalid {
	return &errConfigInvalid{invalidConfig: invalidConfig, wrappedErr: wrappedErr}
}

func (e *errConfigInvalid) Error() string {
	return fmt.Sprintf("%s invalid: %v", e.invalidConfig, e.wrappedErr)
}

func (e *errConfigInvalid) InvalidConfig() string {
	return e.invalidConfig
}

func (e *errConfigInvalid) Unwrap() error {
	return e.wrappedErr
}
