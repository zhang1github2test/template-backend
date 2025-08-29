package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"template-backend/config"
	_ "template-backend/internal/handler"
	"template-backend/internal/middleware"
	"template-backend/internal/router"
	"template-backend/pkg/logger"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func ServerMain() {
	cfg := config.LoadConfig()
	logger.Init()

	logger := logger.Logger()
	defer logger.Sync()

	r := gin.New()
	r.Use(gin.Recovery(), middleware.LoggingMiddleware(logger), middleware.CORSMiddleware(), middleware.JWTMiddleware())

	// 自动注册路由（模块通过 init 注册）
	router.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	logger.Info("Server running", zap.Int("port", cfg.App.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", zap.Error(err))
	}
	logger.Info("Server exiting")
}
