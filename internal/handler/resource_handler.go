package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"template-backend/internal/dto"
	"template-backend/internal/repository"
	"template-backend/internal/router"
	"template-backend/internal/service"
	"template-backend/pkg/utils"
)

type ResourceHandler struct {
	resourceService service.ResourceService
}

func NewResourceHandler(resourceService service.ResourceService) *ResourceHandler {
	return &ResourceHandler{
		resourceService: resourceService,
	}
}

// CreateResource 创建资源
// @Summary 创建资源
// @Description 创建新的资源
// @Tags 资源管理
// @Accept json
// @Produce json
// @Param resource body dto.CreateResourceRequest true "资源信息"
// @Success 200 {object} dto.ResourceResponseDoc
// @Failure 400 {object} dto.ResourceResponseDoc
// @Router /api/resources [post]
func (h *ResourceHandler) CreateResource(c *gin.Context) {
	var req dto.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"message": err.Error(),
		})
		return
	}

	response, err := h.resourceService.CreateResource(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "创建失败",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    response,
	})
}

// GetResource 获取资源详情
// @Summary 获取资源详情
// @Description 根据ID获取资源详情
// @Tags 资源管理
// @Produce json
// @Param id path int true "资源ID"
// @Success 200 {object} dto.ResourceResponseDoc
// @Failure 400 {object} dto.ResourceResponseDoc
// @Router /api/resources/{id} [get]
func (h *ResourceHandler) GetResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"message": "无效的资源ID",
		})
		return
	}

	response, err := h.resourceService.GetResourceByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "获取失败",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    response,
	})
}

// UpdateResource 更新资源
// @Summary 更新资源
// @Description 更新指定ID的资源
// @Tags 资源管理
// @Accept json
// @Produce json
// @Param id path int true "资源ID"
// @Param resource body dto.UpdateResourceRequest true "资源信息"
// @Success 200 {object} dto.ResourceResponseDoc
// @Failure 400 {object} dto.ResourceResponseDoc
// @Router /api/resources/{id} [put]
func (h *ResourceHandler) UpdateResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"message": "无效的资源ID",
		})
		return
	}

	var req dto.UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"message": err.Error(),
		})
		return
	}

	response, err := h.resourceService.UpdateResource(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "更新失败",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    response,
	})
}

// DeleteResource 删除资源
// @Summary 删除资源
// @Description 删除指定ID的资源
// @Tags 资源管理
// @Produce json
// @Param id path int true "资源ID"
// @Success 200 {object} dto.ResourceResponseDoc
// @Failure 400 {object} dto.ResourceResponseDoc
// @Router /api/resources/{id} [delete]
func (h *ResourceHandler) DeleteResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"message": "无效的资源ID",
		})
		return
	}

	err = h.resourceService.DeleteResource(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "删除失败",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// ListResources 查询资源列表
// @Summary 查询资源列表
// @Description 根据条件查询资源列表，支持分页
// @Tags 资源管理
// @Produce json
// @Param id query int false "资源ID"
// @Param resource_name query string false "资源名称"
// @Param permission_code query string false "权限标识码"
// @Param type query string false "资源类型"
// @Param resource_path query string false "资源路径"
// @Param http_method query string false "HTTP方法"
// @Param parent_id query int false "父权限ID"
// @Param status query int false "状态"
// @Param requires_auth query int false "是否需要鉴权"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.ResourcePageResponseDoc
// @Failure 400 {object} dto.ResourcePageResponseDoc
// @Router /api/resources [get]
func (h *ResourceHandler) ListResources(c *gin.Context) {
	var req dto.ResourceQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"message": err.Error(),
		})
		return
	}

	response, err := h.resourceService.ListResources(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "查询失败",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data":    response,
	})
}

func (h *ResourceHandler) ResourcesTree(c *gin.Context) {
	var req dto.ResourceQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"message": err.Error(),
		})
		return
	}

	response, err := h.resourceService.ResourcesTree(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "查询失败",
			"message": err.Error(),
		})
		return
	}
	utils.JSON(c, utils.Success(response))

	//c.JSON(http.StatusOK, gin.H{
	//	"code":    200,
	//	"message": "查询成功",
	//	"data":    response,
	//})
}

func init() {
	// 自动注册路由模块（通过 init 自动调用）
	router.RegisterRouteModule(&ResourceHandler{})
}

func (h *ResourceHandler) Register(rg *gin.RouterGroup, db *gorm.DB) {
	h.resourceService = service.NewResourceService(repository.NewResourceRepository(db))
	resources := rg.Group("/resources")
	{
		resources.POST("", h.CreateResource)
		resources.GET("", h.ListResources)
		resources.GET("/:id", h.GetResource)
		resources.PUT("/:id", h.UpdateResource)
		resources.DELETE("/:id", h.DeleteResource)
		resources.GET("/tree", h.ResourcesTree)
	}

}
