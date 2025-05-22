package errs

import (
	"errors"
)

// 工具函数

// IsErrorCode 检查首个 Code 错误是否为特定错误代码
func IsErrorCode(err error, codes ...code) bool {
	if err == nil {
		return false
	}

	var codeErr WithCodeErr
	if errors.As(err, &codeErr) {
		errCode := codeErr.Code()
		for _, code := range codes {
			if errCode == code {
				return true
			}
		}
	}

	return false
}

func isCodeInRange(code WithCodeErr, start code, end code) bool {
	return code.Code() >= start && code.Code() <= end
}

// ---------------- 错误分类工具函数 ----------------
func isErrorInRange(err error, start code, end code) bool {
	var codeErr WithCodeErr
	if !errors.As(err, &codeErr) {
		return false
	}

	return isCodeInRange(codeErr, start, end)
}

// HasErrorInRange 递归查找整个链路是否在范围内，start == end 具体查找某个 code
func HasErrorInRange(err error, start code, end code) bool {
	if err == nil {
		return false
	}
	var codeErr WithCodeErr

	if errors.As(err, &codeErr) {
		if isCodeInRange(codeErr, start, end) {
			return true
		}

		for _, cause := range codeErr.Unwrap() {
			if HasErrorInRange(cause, start, end) {
				return true
			}
		}
	} else {
		if unwrappedErr := errors.Unwrap(err); unwrappedErr != nil {
			if HasErrorInRange(unwrappedErr, start, end) {
				return true
			}
		}
	}
	return false
}

// IsBadRequest 检查是否为请求相关错误
func IsBadRequest(err error) bool {
	return isErrorInRange(err, 1000, 1999)
}

// IsAuthentication 检查是否为认证相关错误
func IsAuthentication(err error) bool {
	return isErrorInRange(err, 2000, 2999)
}

// IsInternal 检查是否为内部错误
func IsInternal(err error) bool {
	return isErrorInRange(err, 3000, 3999)
}

// IsDatabase 检查是否为数据库相关错误
func IsDatabase(err error) bool {
	return isErrorInRange(err, 4000, 4999)
}

// IsBusiness 检查是否为业务相关错误
func IsBusiness(err error) bool {
	return isErrorInRange(err, 5000, 5999)
}
