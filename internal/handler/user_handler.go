// handler/user_handler.go
package handler

import (
	"net/http"
	"strconv"
	"template-backend/internal/model"
	"template-backend/internal/repository"
	"template-backend/internal/router"
	"template-backend/internal/service"
	"template-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GET /api/users
func (h *UserHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	filters := map[string]interface{}{}
	if username := c.Query("username"); username != "" {
		filters["username"] = username
	}
	if nickname := c.Query("nickname"); nickname != "" {
		filters["nickname"] = nickname
	}
	if email := c.Query("email"); email != "" {
		filters["email"] = email
	}
	if status := c.Query("status"); status != "" {
		if v, err := strconv.Atoi(status); err == nil {
			filters["status"] = v
		}
	}

	users, total, err := h.userService.GetList(page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":     users,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GET /api/users/:id
func (h *UserHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": user})
}

// POST /api/users
func (h *UserHandler) Create(c *gin.Context) {
	var req model.UserCreateInformation
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	var user model.User
	utils.DeepCopyStruct(&user, &req)
	password, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Password 设置不正确"})
		return
	}
	user.Password = password

	if err := h.userService.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	if req.RoleIds != nil {
		if err := h.userService.AssignRoles(user.ID, req.RoleIds); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": req})
}

// PUT /api/users/:id
func (h *UserHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
		return
	}

	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	user.Nickname = req.Nickname
	user.Email = req.Email
	user.Phone = req.Phone
	user.Gender = req.Gender
	user.Status = req.Status

	if err := h.userService.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	if req.RoleIds != nil {
		if err := h.userService.AssignRoles(user.ID, req.RoleIds); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": user})
}

// DELETE /api/users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.userService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *UserHandler) AssignRoles(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		RoleIDs []uint `json:"roleIds" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := h.userService.AssignRoles(uint(userID), req.RoleIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "角色分配失败: " + err.Error()})
		return
	}

	// 获取更新后的用户信息
	userWithRoles, err := h.userService.GetByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": userWithRoles, "message": "角色分配成功"})
}

// GET /api/users/:id/roles - 获取用户的角色
func (h *UserHandler) GetUserRoles(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))

	roles, err := h.userService.GetUserRoles(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户角色失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": roles})
}

func (h *UserHandler) Register(rg *gin.RouterGroup, db *gorm.DB) {
	h.userService = service.NewUserService(repository.NewUserRepository(db))
	users := rg.Group("/users")
	users.GET("", h.GetList)
	users.GET("/:id", h.GetByID)
	users.POST("", h.Create)
	users.PUT("/:id", h.Update)
	users.DELETE("/:id", h.Delete)

	// 新增角色相关路由
	users.POST("/:id/roles", h.AssignRoles)
	users.GET("/:id/roles", h.GetUserRoles)
}

func init() {
	// 自动注册路由模块（通过 init 自动调用）
	router.RegisterRouteModule(&UserHandler{})
}
