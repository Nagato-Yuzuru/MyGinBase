//go:generate stringer -type=errorCode
package errs

type Code int

type WithCodeErr interface {
	error
	Code() Code
}

// 通用基础错误 (1-999)
const (
	ErrUnknown      Code = iota + 1 // 未知错误
	ErrInvalidParam                 // 无效参数
	ErrNotFound                     // 资源未找到
	ErrUnauthorized                 // 未授权（未登录）
	ErrForbidden                    // 禁止访问（无权限）
	ErrInternal                     // 内部服务错误
)

// 请求相关错误 (1000-1999)
const (
	ErrBadRequest           Code = 1000 + iota // 错误的请求格式
	ErrRateLimited                             // 请求频率超限
	ErrTimeout                                 // 请求超时
	ErrPayloadTooLarge                         // 请求内容过大
	ErrUnsupportedMediaType                    // 不支持的媒体类型
)

// 数据操作相关错误 (2000-2999)
const (
	ErrConflict         Code = 2000 + iota // 资源冲突（如唯一键冲突）
	ErrGone                                // 资源不再可用
	ErrDataCorruption                      // 数据损坏或格式错误
	ErrValidationFailed                    // 数据验证失败
)

// 服务状态相关错误 (3000-3999)
const (
	ErrServiceUnavailable Code = 3000 + iota // 服务暂时不可用
	ErrMaintenanceMode                       // 系统维护中
	ErrOverloaded                            // 系统过载
	ErrDependencyFailure                     // 依赖服务失败
)

// 认证与授权相关错误 (4000-4999)
const (
	ErrTokenExpired      Code = 4000 + iota // 令牌过期
	ErrInvalidToken                         // 无效的令牌
	ErrInsufficientScope                    // 权限范围不足
	ErrAccountLocked                        // 账户已锁定
	ErrAccountDisabled                      // 账户已禁用
)

// 数据库相关错误 (5000-5999)
const (
	ErrDBConnection  Code = 5000 + iota // 数据库连接失败
	ErrDBTransaction                    // 事务错误
	ErrDBConstraint                     // 约束冲突
	ErrDBDeadlock                       // 数据库死锁
)
