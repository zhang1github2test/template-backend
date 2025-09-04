package utils

import "github.com/gin-gonic/gin"

// ApiResponse 使用泛型来指定数据类型
type ApiResponse[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data,omitempty"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type PageResult[T any] struct {
	List     []T   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

// Success 使用泛型创建成功的响应
func Success[T any](data T) *ApiResponse[T] {
	return &ApiResponse[T]{
		Success: true,
		Data:    data,
		Message: "ok",
		Code:    200,
	}
}

// Error 创建错误响应（错误响应通常不包含数据）
func Error(message string, code int) *ApiResponse[any] {
	return &ApiResponse[any]{
		Success: false,
		Message: message,
		Code:    code,
	}
}

// JSON 响应函数，支持泛型
func JSON[T any](ctx *gin.Context, resp *ApiResponse[T]) {
	ctx.JSON(resp.Code, resp)
}
