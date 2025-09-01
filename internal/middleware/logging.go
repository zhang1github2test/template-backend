package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"template-backend/internal/global"
	"template-backend/internal/model"
	"template-backend/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ResponseWriter 是一个包装的 ResponseWriter，用于捕获响应内容
type ResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w ResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w ResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Logger     *zap.Logger
	SkipPaths  []string
	MaxBodyLen int
}

// LoggingMiddlewareWithConfig 带配置的日志中间件
func LoggingMiddlewareWithConfig(config LoggingConfig) gin.HandlerFunc {
	skip := make(map[string]bool)
	for _, path := range config.SkipPaths {
		skip[path] = true
	}

	maxBodyLen := config.MaxBodyLen
	if maxBodyLen <= 0 {
		maxBodyLen = 1024 // 默认最大1KB
	}

	return func(c *gin.Context) {
		for _, url := range config.SkipPaths {
			if c.Request.Method == http.MethodOptions || strings.Contains(c.Request.URL.Path, url) {
				logger.Logger().Info("no need to log this url ", zap.String("url", url), zap.String("method", c.Request.Method))
				c.Next()
				return
			}
		}

		start := time.Now()

		// 读取请求体（添加长度限制）
		var requestBody interface{}
		if c.Request.Body != nil && c.Request.Method != http.MethodGet {
			// 限制读取的字节数
			bodyBytes, _ := io.ReadAll(io.LimitReader(c.Request.Body, int64(maxBodyLen)))
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 重新设置 Body

			// 尝试解析为 JSON
			if json.Valid(bodyBytes) {
				var data interface{}
				if err := json.Unmarshal(bodyBytes, &data); err == nil {
					requestBody = data
				} else {
					requestBody = string(bodyBytes)
				}
			} else {
				requestBody = string(bodyBytes)
			}
		}

		// 包装 ResponseWriter 以捕获响应
		blw := &ResponseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(start)

		// 获取客户端IP
		clientIP := c.ClientIP()

		// 获取 User-Agent
		userAgent := c.Request.UserAgent()
		var truncated = false
		// 获取响应内容（同样添加长度限制）
		var responseBody interface{}
		responseBodyBytes := blw.body.Bytes()
		// 限制响应体大小
		if len(responseBodyBytes) > maxBodyLen {
			responseBodyBytes = responseBodyBytes[:maxBodyLen]
			responseBody = string(responseBodyBytes) + "...(truncated)"
			truncated = true
		} else if len(responseBodyBytes) > 0 {
			if json.Valid(responseBodyBytes) {
				var data interface{}
				if err := json.Unmarshal(responseBodyBytes, &data); err == nil {
					responseBody = data
				} else {
					responseBody = string(responseBodyBytes)
				}
			} else {
				responseBody = string(responseBodyBytes)
			}
		}

		// 创建日志对象，符合 model.Log 结构
		logEntry := &model.Log{
			Timestamp:     time.Now(),
			Method:        c.Request.Method,
			Path:          c.Request.URL.Path,
			Query:         c.Request.URL.RawQuery,
			IP:            clientIP,
			UserAgent:     userAgent,
			Status:        c.Writer.Status(),
			Latency:       latency.Milliseconds(),
			Handler:       c.HandlerName(),
			Request:       requestBody,
			Response:      responseBody,
			Errors:        "", // 错误信息将在下面处理
			ContentLength: c.Request.ContentLength,
			Truncated:     truncated,
			CreatedAt:     time.Now(),
		}
		// 处理错误信息
		if len(c.Errors) > 0 {
			errors := make([]string, len(c.Errors))
			for i, err := range c.Errors {
				errors[i] = err.Error()
			}

			logEntry.Errors = strings.Join(errors, "\n")
		}

		select {
		case global.GetLogChan() <- logEntry: // 这里应该发送到通道
		default:
			// 通道满时丢弃日志，避免阻塞
		}

		// 记录日志
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", clientIP),
			zap.String("user-agent", userAgent),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("handler", c.HandlerName()),
		}

		// 只有当请求体不为空时才记录
		if requestBody != nil {
			fields = append(fields, zap.Any("request", requestBody))
		}

		// 只有当响应体不为空时才记录
		if responseBody != nil {
			fields = append(fields, zap.Any("response", responseBody))
		}

		// 记录错误信息
		if len(c.Errors) > 0 {
			fields = append(fields, zap.Strings("errors", c.Errors.Errors()))
		}

		config.Logger.Info("http_request", fields...)
	}
}

// 默认的日志中间件
func EnhancedLoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return LoggingMiddlewareWithConfig(LoggingConfig{
		Logger: logger,
		SkipPaths: []string{
			"/health",
			"/metrics",
			"/system/log",
		},
		MaxBodyLen: 1024 * 10, // 10KB
	})
}
