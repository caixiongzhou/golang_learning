package middleware

import (
	"blog-system/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware 日志中间件
func LoggingMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		
		// 处理请求
		ctx.Next()
		
		// 记录日志
		duration := time.Since(start)
		status := ctx.Writer.Status()
		
		log.Info("HTTP请求: %s %s %d %s",
			ctx.Request.Method,
			ctx.Request.URL.Path,
			status,
			duration,
		)
	}
}