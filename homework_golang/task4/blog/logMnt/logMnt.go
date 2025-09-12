package logMnt

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AppError struct {
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

var (
	ErrDatabaseConnection = &AppError{
		StatusCode: http.StatusInternalServerError,
		Message:    "数据库连接失败",
	}

	ErrUnauthorized = &AppError{
		StatusCode: http.StatusUnauthorized,
		Message:    "用户认证失败",
	}

	ErrForbidden = &AppError{
		StatusCode: http.StatusForbidden,
		Message:    "没有权限执行此操作",
	}

	ErrNotFound = &AppError{
		StatusCode: http.StatusNotFound,
		Message:    "资源不存在",
	}

	ErrBadRequest = &AppError{
		StatusCode: http.StatusBadRequest,
		Message:    "请求参数错误",
	}

	ErrInternalServerError = &AppError{
		StatusCode: http.StatusInternalServerError,
		Message:    "服务器内部错误",
	}
)

// 初始化zap日志配置
func InitZapLogger() (*zap.Logger, error) {
	// 开发环境配置
	config := zap.NewDevelopmentConfig()

	// 设置日志级别
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	// 设置时间格式
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	return config.Build()
}

// 日志中间件
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(startTime)

		// 记录日志
		zap.L().Info("请求处理完成",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.String("ip", c.ClientIP()),
		)
	}
}

// 错误处理中间件
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 处理请求
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			var appErr *AppError

			// 判断错误类型
			if errors.As(err, &appErr) {
				// 已知应用错误
				zap.L().Error("应用错误",
					zap.String("message", appErr.Message),
					zap.Int("status", appErr.StatusCode),
					zap.Error(appErr.Err),
					zap.String("path", c.Request.URL.Path),
				)
				c.JSON(appErr.StatusCode, gin.H{
					"success": false,
					"error":   appErr.Message,
				})
			} else {
				// 未知错误
				zap.L().Error("未知错误",
					zap.Error(err),
					zap.String("path", c.Request.URL.Path),
				)
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   "服务器内部错误",
				})
			}
			c.Abort()
		}

	}
}
