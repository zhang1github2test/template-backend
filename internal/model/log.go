package model

import (
	"gorm.io/gorm"
	"time"
)

// Log HTTP请求日志结构体
type Log struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Timestamp     time.Time      `json:"timestamp"`
	Method        string         `json:"method"`
	Path          string         `json:"path"`
	Query         string         `json:"query"`
	IP            string         `json:"ip"`
	UserAgent     string         `json:"userAgent"`
	Status        int            `json:"status"`
	Latency       int64          `json:"latency"`
	Handler       string         `json:"handler"`
	Request       interface{}    `json:"request" gorm:"serializer:json"`
	Response      interface{}    `json:"response" gorm:"serializer:json"`
	Errors        string         `json:"errors"`
	ContentLength int64          `json:"contentLength"`
	Truncated     bool           `json:"truncated"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
