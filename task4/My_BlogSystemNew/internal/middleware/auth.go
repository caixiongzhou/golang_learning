package middleware

import (
	"blog-system/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
// 该中间件用于验证HTTP请求中的JWT令牌，保护需要认证的接口
func AuthMiddleware(jwtUtil *utils.JWTUtil) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 步骤1：从请求头中获取Authorization字段
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			// 如果没有提供认证令牌，返回401未授权错误
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "缺少认证令牌",
				"data":    nil,
			})
			// 中止后续处理链
			ctx.Abort()
			return
		}

		// 步骤2：检查Bearer token格式是否正确
		// 正确的格式: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// 令牌格式不符合Bearer规范，返回401错误
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "令牌格式错误，应为: Bearer <token>",
				"data":    nil,
			})
			ctx.Abort()
			return
		}

		// 步骤3：提取并验证JWT令牌
		token := parts[1] // 获取实际的令牌字符串
		claims, err := jwtUtil.ValidateToken(token)
		if err != nil {
			// 令牌验证失败（可能过期、签名错误或被篡改）
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "令牌无效或已过期",
				"data":    nil,
			})
			ctx.Abort()
			return
		}

		// 步骤4：令牌验证成功，将用户信息存储到Gin上下文
		// 这样后续的控制器就可以从ctx中获取用户信息
		ctx.Set("userID", claims.UserID)     // 设置用户ID到上下文
		ctx.Set("username", claims.Username) // 设置用户名到上下文

		// 步骤5：认证通过，继续执行后续的中间件和请求处理函数
		ctx.Next()
	}
}
