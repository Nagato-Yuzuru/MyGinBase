package errs

import (
	"fmt"
	"strings"
)

type baseError struct {
	code           // 错误代码
	causes []error // 原始错误
}

func (e *baseError) Is(target error) bool {

	codeErr, ok := target.(WithCodeErr)

	if !ok {
		return false
	}
	if codeErr.Code() == e.code {
		return true
	}
	return false
}

func (e *baseError) Error() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("[%d] %s", e.code, e.code))

	if len(e.causes) == 0 {
		return b.String()
	}

	for _, cause := range e.causes {
		b.WriteString(", ")
		b.WriteString(cause.Error())
	}

	return b.String()
}

func (e *baseError) Code() code {
	return e.code
}

func (e *baseError) Unwrap() []error {
	return e.causes
}

func newError(code code, causes ...error) *baseError {

	return &baseError{
		code:   code,
		causes: causes,
	}
}

func WrapCodeError(code code, causes ...error) WithCodeErr {
	return newError(code, causes...)
}
