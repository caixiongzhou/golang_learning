package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 成功响应
func Success(ctx *gin.Context, status int, message string, data interface{}) {
	ctx.JSON(status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, Response{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

// ValidationError 参数验证错误
func ValidationError(ctx *gin.Context, errors map[string]string) {
	ctx.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: "参数验证失败",
		Data:    errors,
	})
}

// 说明：统一的响应结构与工具方法，方便控制器返回一致格式的 JSON 响应。
// - Success 为是否成功，Message 为简短提示，Data 放置业务数据或错误详情。