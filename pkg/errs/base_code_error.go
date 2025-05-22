//go:generate stringer -type=code
package errs

import (
	"fmt"
)

type code int

type BaseChainError interface {
	error
	Unwrap() []error
	Is(error) bool
}

type WithCodeErr interface {
	fmt.Stringer
	BaseChainError
	Code() code
}

// 通用基础错误 (1-999)
const (
	ErrUnknown code = iota + 1 // 未知错误
	// 注意：原有的 ErrInvalidParam, ErrUnauthorized, ErrForbidden, ErrNotFound, ErrInternal 已移至更具体的分类
)

// 请求处理相关错误 (1000-1999)
const (
	ErrBadRequest           code = 1000 + iota // 错误的请求格式
	ErrInvalidParam                            // 无效参数 (从通用错误移入)
	ErrNotFound                                // 资源未找到 (从通用错误移入)
	ErrConflict                                // 资源冲突（如唯一键冲突, 从数据操作移入）
	ErrGone                                    // 资源不再可用 (从数据操作移入)
	ErrValidationFailed                        // 数据验证失败 (从数据操作移入)
	ErrRateLimited                             // 请求频率超限
	ErrTimeout                                 // 请求超时 (可区分为客户端请求超时或服务端处理超时)
	ErrPayloadTooLarge                         // 请求内容过大
	ErrUnsupportedMediaType                    // 不支持的媒体类型
)

// 认证与授权相关错误 (2000-2999)
const (
	ErrUnauthorized      code = 2000 + iota // 未授权（通常指未认证或认证失败, 从通用错误移入）
	ErrForbidden                            // 禁止访问（已认证但无权限, 从通用错误移入）
	ErrTokenExpired                         // 令牌过期
	ErrInvalidToken                         // 无效的令牌
	ErrInsufficientScope                    // 权限范围不足
	ErrAccountLocked                        // 账户已锁定
	ErrAccountDisabled                      // 账户已禁用
)

// 服务端状态与内部错误 (3000-3999)
const (
	ErrInternalServer     code = 3000 + iota // 内部服务错误 (合并了原 ErrInternal 和 ErrInternalServerError)
	ErrServiceUnavailable                    // 服务暂时不可用
	ErrMaintenanceMode                       // 系统维护中
	ErrOverloaded                            // 系统过载
	ErrDependencyFailure                     // 依赖服务失败
	// ErrTimeoutServiceUnavailable 可以被 ErrTimeout (如果明确为服务端) 或 ErrDependencyFailure 覆盖，或保留并更名
	ErrNotImplemented // 功能未实现 (从原内部错误移入)
	ErrNilPointer     // 空指针异常 (从原运行时错误移入)
	ErrEnvironmentConfig
)

// 数据存储相关错误 (4000-4999)
const (
	ErrDBConnection   code = 4000 + iota // 数据库连接失败
	ErrDBTransaction                     // 事务错误
	ErrDBConstraint                      // 约束冲突
	ErrDBDeadlock                        // 数据库死锁
	ErrDataCorruption                    // 数据损坏或格式错误 (从原数据操作移入)
)

// 业务错误 （5000-5999）
const ()
