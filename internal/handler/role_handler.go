// handler/role_handler.go
package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"template-backend/internal/model"
	"template-backend/internal/repository"
	"template-backend/internal/router"
	"template-backend/internal/service"
	"template-backend/pkg/utils"
)

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

// GET /api/roles
func (h *RoleHandler) GetRoleList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	filters := map[string]interface{}{}
	if roleName := c.Query("roleName"); roleName != "" {
		filters["roleName"] = roleName
	}
	if roleCode := c.Query("roleCode"); roleCode != "" {
		filters["roleCode"] = roleCode
	}
	if status := c.Query("status"); status != "" {
		if v, err := strconv.Atoi(status); err == nil {
			filters["status"] = v
		}
	}

	roles, total, err := h.roleService.GetList(page, size, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	utils.JSON(c, utils.Success(gin.H{
		"list":  roles,
		"total": total,
		"page":  page,
		"size":  size,
	}))

}

// POST /api/roles
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req model.Role
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := h.roleService.Create(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": req})
}

// PUT /api/roles/:id
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	role, err := h.roleService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
		return
	}

	var req model.Role
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	role.RoleName = req.RoleName
	role.RoleCode = req.RoleCode
	role.RoleDesc = req.RoleDesc
	role.Status = req.Status

	if err := h.roleService.Update(role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": role})
}

// DELETE /api/roles/:id
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.roleService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

// DELETE /api/roles/batch
func (h *RoleHandler) BatchDeleteRoles(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := h.roleService.BatchDelete(req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "批量删除成功"})
}

// GET /api/roles/:id/permissions
func (h *RoleHandler) GetRolePermissions(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	permissions, err := h.roleService.GetPermissions(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": permissions})
}

// PUT /api/roles/:id/permissions
func (h *RoleHandler) UpdateRolePermissions(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		PermissionIds []uint `json:"permissionIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := h.roleService.UpdatePermissions(uint(id), req.PermissionIds); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "权限更新成功"})
}

func (h *RoleHandler) Register(rg *gin.RouterGroup, db *gorm.DB) {
	h.roleService = service.NewRoleService(repository.NewRoleRepository(db))
	roles := rg.Group("/roles")
	{
		roles.GET("", h.GetRoleList)
		roles.POST("", h.CreateRole)
		roles.PUT("/:id", h.UpdateRole)
		roles.DELETE("/:id", h.DeleteRole)
		roles.DELETE("/batch", h.BatchDeleteRoles)
		roles.GET("/:id/permissions", h.GetRolePermissions)
		roles.PUT("/:id/permissions", h.UpdateRolePermissions)
	}
}

func init() {
	// 自动注册路由模块（通过 init 自动调用）
	router.RegisterRouteModule(&RoleHandler{})
}
