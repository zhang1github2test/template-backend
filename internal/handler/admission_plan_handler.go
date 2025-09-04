package handler

import (
	"net/http"
	"strconv"
	"template-backend/internal/repository"
	"template-backend/internal/router"

	"gorm.io/gorm"

	"template-backend/internal/model"
	"template-backend/internal/service"
	"template-backend/pkg/logger"
	"template-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ListRequest 查询参数结构体
type ListRequest struct {
	Page         int    `form:"page" json:"page"`
	PageSize     int    `form:"page_size" json:"page_size"`
	SchoolName   string `form:"school_name" json:"school_name"`
	Year         int    `form:"year" json:"year"`
	DistrictType string `form:"district_type" json:"district_type"`
}

type AdmissionPlanHandler struct {
	service *service.AdmissionPlanService
}

func NewAdmissionPlanHandler(svc *service.AdmissionPlanService) *AdmissionPlanHandler {
	return &AdmissionPlanHandler{service: svc}
}

// @Summary 获取招生计划列表
// @Description 分页获取招生计划信息，支持条件筛选
// @Tags 招生计划
// @Accept json
// @Produce json
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页数量，默认10"
// @Param school_name query string false "学校名称"
// @Param year query int false "年份"
// @Param district_type query string false "区属类型"
// @Success 200 {object} utils.ApiResponse
// @Failure 400 {object} utils.ApiResponse
// @Router /api/plans [get]
func (h *AdmissionPlanHandler) List(c *gin.Context) {
	var req ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Logger().Error("List 参数绑定失败", zap.Error(err))
		utils.JSON(c, utils.Error("参数错误", http.StatusBadRequest))
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	logger.Logger().Info("List 入参", zap.Any("request", req))

	filters := make(map[string]interface{})
	if req.SchoolName != "" {
		filters["school_name"] = req.SchoolName
	}
	if req.Year != 0 {
		filters["year"] = req.Year
	}
	if req.DistrictType != "" {
		filters["district_type"] = req.DistrictType
	}

	plans, total, err := h.service.List(req.Page, req.PageSize, filters)
	if err != nil {
		logger.Logger().Error("List 查询失败", zap.Error(err))
		utils.JSON(c, utils.Error("查询失败", http.StatusInternalServerError))
		return
	}

	logger.Logger().Info("List 查询成功", zap.Int64("total", total))
	//utils.JSON(c, utils.Success(gin.H{
	//	"list":      plans,
	//	"total":     total,
	//	"page":      req.Page,
	//	"page_size": req.PageSize,
	//}))
	utils.JSON(c, utils.Success(utils.PageResult[model.HighSchoolAdmissionPlan]{
		List:     plans,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}))
}

// @Summary 获取单个招生计划详情
// @Description 根据ID获取招生计划详情
// @Tags 招生计划
// @Accept json
// @Produce json
// @Param id path int true "招生计划ID"
// @Success 200 {object} utils.ApiResponse
// @Failure 400 {object} utils.ApiResponse
// @Router /api/plans/{id} [get]
func (h *AdmissionPlanHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Logger().Error("GetByID ID参数错误", zap.Error(err))
		utils.JSON(c, utils.Error("参数错误", http.StatusBadRequest))
		return
	}

	logger.Logger().Info("GetByID 入参", zap.Int("id", id))

	plan, err := h.service.GetByID(id)
	if err != nil {
		logger.Logger().Error("GetByID 查询失败", zap.Error(err), zap.Int("id", id))
		utils.JSON(c, utils.Error("未找到数据", http.StatusNotFound))
		return
	}

	logger.Logger().Info("GetByID 查询成功", zap.Int("id", id))
	utils.JSON(c, utils.Success(plan))
}

// @Summary 创建招生计划
// @Description 新增一条招生计划记录
// @Tags 招生计划
// @Accept json
// @Produce json
// @Param plan body model.HighSchoolAdmissionPlan true "招生计划对象"
// @Success 200 {object} utils.ApiResponse
// @Failure 400 {object} utils.ApiResponse
// @Router /api/plans [post]
func (h *AdmissionPlanHandler) Create(c *gin.Context) {
	var plan model.HighSchoolAdmissionPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		logger.Logger().Error("Create 参数绑定失败", zap.Error(err))
		utils.JSON(c, utils.Error("参数错误", http.StatusBadRequest))
		return
	}

	logger.Logger().Info("Create 入参", zap.Any("plan", plan))

	if err := h.service.Create(&plan); err != nil {
		logger.Logger().Error("Create 创建失败", zap.Error(err), zap.Any("plan", plan))
		utils.JSON(c, utils.Error("创建失败", http.StatusInternalServerError))
		return
	}

	logger.Logger().Info("Create 创建成功", zap.Int("id", plan.ID))
	utils.JSON(c, utils.Success(plan))
}

// @Summary 更新招生计划
// @Description 根据ID更新招生计划信息
// @Tags 招生计划
// @Accept json
// @Produce json
// @Param id path int true "招生计划ID"
// @Param plan body model.HighSchoolAdmissionPlan true "更新内容"
// @Success 200 {object} utils.ApiResponse
// @Failure 400 {object} utils.ApiResponse
// @Router /api/plans/{id} [put]
func (h *AdmissionPlanHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Logger().Error("Update ID参数错误", zap.Error(err))
		utils.JSON(c, utils.Error("参数错误", http.StatusBadRequest))
		return
	}

	var plan model.HighSchoolAdmissionPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		logger.Logger().Error("Update 参数绑定失败", zap.Error(err))
		utils.JSON(c, utils.Error("参数错误", http.StatusBadRequest))
		return
	}

	logger.Logger().Info("Update 入参", zap.Int("id", id), zap.Any("plan", plan))

	if err := h.service.Update(id, &plan); err != nil {
		logger.Logger().Error("Update 更新失败", zap.Error(err), zap.Int("id", id))
		utils.JSON(c, utils.Error("更新失败", http.StatusInternalServerError))
		return
	}

	logger.Logger().Info("Update 更新成功", zap.Int("id", id))
	utils.JSON(c, utils.Success(plan))
}

// @Summary 删除招生计划
// @Description 根据ID删除招生计划
// @Tags 招生计划
// @Accept json
// @Produce json
// @Param id path int true "招生计划ID"
// @Success 200 {object} utils.ApiResponse
// @Failure 400 {object} utils.ApiResponse
// @Router /api/plans/{id} [delete]
func (h *AdmissionPlanHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Logger().Error("Delete ID参数错误", zap.Error(err))
		utils.JSON(c, utils.Error("参数错误", http.StatusBadRequest))
		return
	}

	logger.Logger().Info("Delete 入参", zap.Int("id", id))

	if err := h.service.Delete(id); err != nil {
		logger.Logger().Error("Delete 删除失败", zap.Error(err), zap.Int("id", id))
		utils.JSON(c, utils.Error("删除失败", http.StatusInternalServerError))
		return
	}

	logger.Logger().Info("Delete 删除成功", zap.Int("id", id))
	utils.JSON(c, utils.Success(""))
}

func init() {
	// 自动注册路由模块（通过 init 自动调用）
	router.RegisterRouteModule(&AdmissionPlanHandler{})
}

func (h *AdmissionPlanHandler) Register(rg *gin.RouterGroup, db *gorm.DB) {
	h.service = service.NewAdmissionPlanService(repository.NewAdmissionPlanRepo(db))
	api := rg.Group("/plans")
	{
		api.GET("", h.List)
		api.GET("/:id", h.GetByID)
		api.POST("", h.Create)
		api.PUT("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
	}

}
