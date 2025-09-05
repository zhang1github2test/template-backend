package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"template-backend/internal/dto"
	"template-backend/internal/model"
	"template-backend/internal/repository"
	"template-backend/internal/router"
	"template-backend/internal/service"
	"template-backend/pkg/logger"
	"template-backend/pkg/utils"
)

type SchoolAdmissionHandler struct {
	svc service.SchoolAdmissionService
}

func init() {
	// 注册路由模块
	// ⚠️ 注意：router.RegisterRouteModule 需要你在 router/router.go 里实现
	// 类似 AuthHandler 那样
	router.RegisterRouteModule(&SchoolAdmissionHandler{})
}

// Register 注册路由
func (h *SchoolAdmissionHandler) Register(rg *gin.RouterGroup, db *gorm.DB) {
	h.svc = service.NewSchoolAdmissionService(repository.NewSchoolAdmissionRepository(db))
	school := rg.Group("/school-admission")
	{
		school.POST("", h.Create)
		school.GET("/:id", h.GetByID)
		school.GET("", h.List)
		school.PUT("/:id", h.Update)
		school.DELETE("/:id", h.Delete)
	}
}

// Create godoc
// @Summary 新增中考录取信息
// @Tags 中考录取线
// @Accept json
// @Produce json
// @Param data body model.SchoolAdmissionInfo true "学校招生信息"
// @Success 200 {object} dto.SchoolAdmissionInfoResponseDoc
// @Failure 400 {object} dto.SchoolAdmissionInfoResponseDoc
// @Router /school-admission [post]
func (h *SchoolAdmissionHandler) Create(c *gin.Context) {
	var req model.SchoolAdmissionInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON(c, utils.Error(err.Error(), http.StatusBadRequest))
		return
	}
	if err := h.svc.Create(&req); err != nil {
		utils.JSON(c, utils.Error(err.Error(), http.StatusInternalServerError))
		return
	}
	utils.JSON(c, utils.Success(req))
}

// GetByID godoc
// @Summary 获取中考录取信息
// @Tags 中考录取线
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} dto.SchoolAdmissionInfoResponseDoc
// @Failure 400 {object} dto.SchoolAdmissionInfoResponseDoc
// @Router /school-admission/{id} [get]
func (h *SchoolAdmissionHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	info, err := h.svc.GetByID(id)
	if err != nil {
		utils.JSON(c, utils.Error(err.Error(), http.StatusNotFound))
		return
	}
	utils.JSON(c, utils.Success(info))
}

// List 获取分页列表
// @Summary 获取中考录取信息列表
// @Description 分页 + 搜索（按学校名称、年份）
// @Tags 中考录取线
// @Accept json
// @Produce json
// @Param page query int true "页码"
// @Param pageSize query int true "每页数量"
// @Param schoolName query string false "学校名称（模糊查询）"
// @Param year query int false "年份"
// @Success 200 {object} dto.SchoolAdmissionInfoPageResponseDoc
// @Failure 400 {object} dto.SchoolAdmissionInfoPageResponseDoc
// @Router /api/school-admission [get]
func (h *SchoolAdmissionHandler) List(c *gin.Context) {
	var req dto.SchoolAdmissionQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.JSON(c, utils.Error(err.Error(), http.StatusBadRequest))
		return
	}
	if req.Page <= 0 || req.PageSize <= 0 {
		logger.Logger().Info("invalid params", zap.Int("page", req.Page), zap.Int("pageSize", req.PageSize))
		utils.JSON(c, utils.Error("invalid params", http.StatusBadRequest))
		return
	}

	list, total, err := h.svc.List(&req)
	if err != nil {
		utils.JSON(c, utils.Error(err.Error(), http.StatusInternalServerError))
		return
	}

	utils.JSON(c, utils.Success(gin.H{
		"list":     list,
		"total":    total,
		"page":     req.Page,
		"pageSize": req.PageSize,
	}))
}

// Update godoc
// @Summary 更新学中考录取信息
// @Tags 中考录取线
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param data body model.SchoolAdmissionInfo true "学校招生信息"
// @Success 200 {object} dto.SchoolAdmissionInfoResponseDoc
// @Failure 400 {object} dto.SchoolAdmissionInfoResponseDoc
// @Router /school-admission/{id} [put]
func (h *SchoolAdmissionHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req model.SchoolAdmissionInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON(c, utils.Error(err.Error(), http.StatusBadRequest))
		return
	}
	req.ID = id
	if err := h.svc.Update(&req); err != nil {
		utils.JSON(c, utils.Error(err.Error(), http.StatusInternalServerError))
		return
	}
	utils.JSON(c, utils.Success(req))
}

// Delete godoc
// @Summary 删除中考录取信息
// @Tags 中考录取线
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} dto.SchoolAdmissionInfoResponseDoc
// @Failure 400 {object} dto.SchoolAdmissionInfoResponseDoc
// @Router /school-admission/{id} [delete]
func (h *SchoolAdmissionHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	admissionInfo, err := h.svc.GetByID(id)
	if err != nil {
		utils.JSON(c, utils.Error(err.Error(), http.StatusNotFound))
		return
	}
	if err := h.svc.Delete(id); err != nil {
		utils.JSON(c, utils.Error(err.Error(), http.StatusInternalServerError))
		return
	}
	utils.JSON(c, utils.Success(admissionInfo))
}
