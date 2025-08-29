package handler

import (
	"net/http"
	"time"

	"template-backend/internal/model"
	"template-backend/internal/router"
	"template-backend/internal/service"
	"template-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct{}

func init() {
	// 自动注册路由模块（通过 init 自动调用）
	router.RegisterRouteModule(&AuthHandler{})
}

func (h *AuthHandler) Register(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/login", h.Login)
	// 需要鉴权的接口使用 JWT 中间件
	auth.POST("/logout", h.Logout)
	auth.GET("/user", h.User)
	auth.POST("/refresh", h.Refresh)
	auth.POST("/change-password", h.ChangePassword)
}

// Login 登录
func (h *AuthHandler) Login(c *gin.Context) {
	var form model.LoginForm
	if err := c.ShouldBindJSON(&form); err != nil {
		utils.JSON(c, utils.Error("invalid params", http.StatusBadRequest))
		return
	}

	resp, err := service.Login(form.Username, form.Password)
	if err != nil {
		utils.JSON(c, utils.Error("internal error", http.StatusInternalServerError))
		return
	}
	if resp == nil {
		utils.JSON(c, utils.Error("username or password incorrect", http.StatusUnauthorized))
		return
	}

	utils.JSON(c, utils.Success(resp))
}

// Logout 登出（演示：客户端删除 token 即可；这里返回成功）
func (h *AuthHandler) Logout(c *gin.Context) {
	utils.JSON(c, utils.Success(gin.H{"ok": true}))
}

// User 获取当前用户信息（从 token 中取 username）
func (h *AuthHandler) User(c *gin.Context) {
	// middleware 已验证 token，有需要可把解析结果放到 context
	usernameIfc, exists := c.Get("username")
	if !exists {
		utils.JSON(c, utils.Error("not authenticated", http.StatusUnauthorized))
		return
	}
	username, _ := usernameIfc.(string)
	user := service.GetUserInfo(username)
	utils.JSON(c, utils.Success(user))
}

// Refresh 刷新 token：解析旧 token 并返回新的 token（简单实现）
func (h *AuthHandler) Refresh(c *gin.Context) {
	type RefreshReq struct {
		Token string `json:"token" binding:"required"`
	}
	var req RefreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON(c, utils.Error("invalid params", http.StatusBadRequest))
		return
	}
	newToken, err := service.RefreshToken(req.Token)
	if err != nil {
		utils.JSON(c, utils.Error("refresh failed: "+err.Error(), http.StatusUnauthorized))
		return
	}
	utils.JSON(c, utils.Success(gin.H{"token": newToken}))
}

// ChangePassword 修改密码（模拟校验旧密码）
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var form model.ChangePasswordForm
	if err := c.ShouldBindJSON(&form); err != nil {
		utils.JSON(c, utils.Error("invalid params", http.StatusBadRequest))
		return
	}

	// TODO: 实际场景需要从 token 中获取用户并校验旧密码（bcrypt）
	usernameIfc, exists := c.Get("username")
	if !exists {
		utils.JSON(c, utils.Error("not authenticated", http.StatusUnauthorized))
		return
	}
	username := usernameIfc.(string)

	ok := service.ChangePassword(username, form.OldPassword, form.NewPassword)
	if !ok {
		utils.JSON(c, utils.Error("old password incorrect", http.StatusBadRequest))
		return
	}
	utils.JSON(c, utils.Success(gin.H{"ok": true, "updatedAt": time.Now().Format(time.RFC3339)}))
}
