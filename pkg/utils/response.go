package utils

import "github.com/gin-gonic/gin"

type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
}

func Success(data interface{}) *ApiResponse {
	return &ApiResponse{
		Success: true,
		Data:    data,
		Message: "ok",
		Code:    200,
	}
}

func Error(message string, code int) *ApiResponse {
	return &ApiResponse{
		Success: false,
		Message: message,
		Code:    code,
	}
}

func JSON(ctx *gin.Context, resp *ApiResponse) {
	ctx.JSON(resp.Code, resp)
}
