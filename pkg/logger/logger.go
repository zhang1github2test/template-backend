package logger

import (
	"go.uber.org/zap"
	"sync"
	"template-backend/config"
)

var (
	log  *zap.Logger
	once sync.Once
)

// Init 初始化全局 logger
func Init() {
	once.Do(func() {
		var err error
		appConfig := config.GetConfig()
		if appConfig.App.Env == "dev" {
			log, err = zap.NewDevelopment()
		} else {
			log, err = zap.NewProduction()
		}
		if err != nil {
			panic(err)
		}
	})
}

// Logger 获取全局 logger
func Logger() *zap.Logger {
	if log == nil {
		panic("logger 未初始化，请先调用 logger.Init()")
	}
	return log
}
