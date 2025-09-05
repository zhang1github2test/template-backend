// internal/handler/auth_handler.go
package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"template-backend/internal/repository"
	"template-backend/internal/router"
	"template-backend/internal/service"
	"template-backend/pkg/utils"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	resp, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error()})
		return
	}
	utils.JSON(c, utils.Success(resp))
}

// GetUserInfo 获取用户信息
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	// 从token中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}

	userInfo, err := h.authService.GetUserInfo(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	utils.JSON(c, utils.Success(userInfo))
}

// RefreshToken 刷新token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	newToken, err := h.authService.RefreshToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"token": newToken}})

	utils.JSON(c, utils.Success(gin.H{"token": newToken}))
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	// 从token中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var req struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	err := h.authService.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	utils.JSON(c, utils.Success(gin.H{"ok": true}))
}

// Logout 登出（演示：客户端删除 token 即可；这里返回成功）
func (h *AuthHandler) Logout(c *gin.Context) {
	utils.JSON(c, utils.Success(gin.H{"ok": true}))
}

func init() {
	// 自动注册路由模块（通过 init 自动调用）
	router.RegisterRouteModule(&AuthHandler{})
}

func (h *AuthHandler) Register(rg *gin.RouterGroup, db *gorm.DB) {
	h.authService = service.NewAuthService(repository.NewUserRepository(db), repository.NewRoleRepository(db))
	auth := rg.Group("/auth")
	auth.POST("/login", h.Login)
	// 需要鉴权的接口使用 JWT 中间件
	auth.POST("/logout", h.Logout)
	auth.GET("/user", h.GetUserInfo)
	auth.POST("/refresh", h.RefreshToken)
	auth.POST("/change-password", h.ChangePassword)
}
