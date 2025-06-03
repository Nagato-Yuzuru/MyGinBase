package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"terraqt.io/colas/bedrock-go/pkg/errs"
	"terraqt.io/colas/bedrock-go/pkg/logger"
)

const UnifiedResponseKey = "unified_response_data"

func Success(c *gin.Context, data any) {
	c.Set(UnifiedResponseKey, data)
}

func Error(c *gin.Context, err error) {
	_ = c.Error(err)
	c.Set(UnifiedResponseKey, nil)
}

type UnifiedResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    any    `json:"data"`
	Warning bool   `json:"warning"`
}

// NewUnifiedResponse creates a UnifiedResponse based on the given error and data, returning HTTP status code and response.
func NewUnifiedResponse(err error, data any) (int, *UnifiedResponse) {
	if err == nil {
		return http.StatusOK, &UnifiedResponse{
			Code:    http.StatusOK,
			Msg:     "",
			Data:    data,
			Warning: false,
		}
	}

	var (
		code int
		msg  string
	)

	var codedErr errs.CodedError
	if errors.As(err, &codedErr) {
		code = errs.GetHttpCodeByError(codedErr)
		msg = errs.Helper(codedErr)
	} else {
		code = http.StatusInternalServerError
		msg = "inner error"
	}

	return code, &UnifiedResponse{
		Code:    code,
		Msg:     msg,
		Data:    nil,
		Warning: true,
	}

}

func ResponseNormalizer(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 如果响应已被写入(如文件下载)，则直接跳过
		if c.Writer.Written() {
			return
		}

		var code int
		var body *UnifiedResponse

		lastErr := c.Errors.Last()
		if lastErr != nil {
			var codedErr errs.CodedError
			if errors.As(lastErr.Err, &codedErr) {
				code, body = NewUnifiedResponse(codedErr, nil)
			} else {
				// 未知的、未包装的错误
				log.Warn(c, "Unknown error occurred: %s\n", zap.Error(lastErr)) // 记录详细日志
				code, body = NewUnifiedResponse(lastErr.Err, nil)
			}
		} else {
			data, _ := c.Get(UnifiedResponseKey)
			// 如果 handler 既没设置数据也没报错，就返回一个默认的成功响应
			code, body = NewUnifiedResponse(nil, data)
		}

		c.JSON(code, body)
	}
}
