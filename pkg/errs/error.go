package errs

import (
	"fmt"
	"strings"
)

type codedError struct {
	code           // 错误代码
	causes []error // 原始错误
}

func (e *codedError) Is(target error) bool {

	codeErr, ok := target.(CodedError)

	if !ok {
		return false
	}
	if codeErr.Code() == e.code {
		return true
	}
	return false
}

func (e *codedError) Error() string {
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

func (e *codedError) Code() code {
	return e.code
}

func (e *codedError) Unwrap() []error {
	return e.causes
}

func newError(code code, causes ...error) *codedError {

	return &codedError{
		code:   code,
		causes: causes,
	}
}

func WrapCodeError(code code, causes ...error) CodedError {
	if len(causes) == 1 {
		if cause, ok := causes[0].(CodedError); ok && cause.Code() == code {
			return cause
		}
	}

	return newError(code, causes...)
}
