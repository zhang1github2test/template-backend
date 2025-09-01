package global

import "template-backend/internal/model"

// 全局日志通道
var logChan chan *model.Log

// 初始化日志通道
func InitLogChan(bufferSize int) {
	logChan = make(chan *model.Log, bufferSize)
}

// GetLogChan 获取日志通道
func GetLogChan() chan *model.Log {
	return logChan
}

// CloseLogChan 关闭日志通道
func CloseLogChan() {
	if logChan != nil {
		close(logChan)
	}
}
