package handler

import (
	"net/http"
	"strconv"
	"template-backend/internal/model"
	"template-backend/internal/repository"
	"template-backend/internal/router"
	"template-backend/internal/service"
	"template-backend/pkg/logger"
	"template-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ConfigHandler struct {
	service service.ConfigService
}

func NewConfigHandler(s service.ConfigService) *ConfigHandler {
	return &ConfigHandler{service: s}
}

// GET /api/system/config/list
func (h *ConfigHandler) GetConfigList(c *gin.Context) {
	configKey := c.Query("configKey")
	configName := c.Query("configName")
	pageNumStr := c.DefaultQuery("pageNum", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	pageNum, _ := strconv.Atoi(pageNumStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	params := map[string]interface{}{
		"configKey":  configKey,
		"configName": configName,
		"pageNum":    pageNum,
		"pageSize":   pageSize,
	}

	logger.Logger().Info("GetConfigList 入参",
		zap.String("configKey", configKey),
		zap.String("configName", configName),
		zap.Int("pageNum", pageNum),
		zap.Int("pageSize", pageSize),
	)
	configs, total, err := h.service.GetList(params)
	if err != nil {
		logger.Logger().Error("GetConfigList 失败", zap.Error(err))
		utils.JSON(c, utils.Error(err.Error(), http.StatusInternalServerError))
		return
	}

	data := map[string]any{
		"rows":  configs,
		"total": total,
	}
	logger.Logger().Info("GetConfigList 出参", zap.Any("data", data))
	utils.JSON(c, utils.Success(data))
}

// GET /api/system/config/:id
func (h *ConfigHandler) GetConfigById(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	logger.Logger().Info("GetConfigById 入参", zap.Int64("id", id))

	config, err := h.service.GetByID(id)
	if err != nil {
		logger.Logger().Error("GetConfigById 失败", zap.Error(err))
		utils.JSON(c, utils.Error(err.Error(), http.StatusInternalServerError))
		return
	}

	logger.Logger().Info("GetConfigById 出参", zap.Any("data", config))
	utils.JSON(c, utils.Success(config))
}

// POST /api/system/config
func (h *ConfigHandler) AddConfig(c *gin.Context) {
	var config model.Config
	if err := c.ShouldBindJSON(&config); err != nil {
		utils.JSON(c, utils.Error(err.Error(), http.StatusInternalServerError))
		return
	}

	logger.Logger().Info("AddConfig 入参", zap.Any("config", config))

	if err := h.service.Create(&config); err != nil {
		logger.Logger().Error("AddConfig 失败", zap.Error(err))
		utils.JSON(c, utils.Error(err.Error(), http.StatusInternalServerError))
		return
	}

	logger.Logger().Info("AddConfig 出参", zap.Any("data", config))
	utils.JSON(c, utils.Success(config))
}

// PUT /api/system/config
func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	var config model.Config
	if err := c.ShouldBindJSON(&config); err != nil {
		utils.JSON(c, utils.Error(err.Error(), http.StatusInternalServerError))
		return
	}

	logger.Logger().Info("UpdateConfig 入参", zap.Any("config", config))

	if err := h.service.Update(&config); err != nil {
		logger.Logger().Error("UpdateConfig 失败", zap.Error(err))
		utils.JSON(c, utils.Error(err.Error(), http.StatusInternalServerError))
		return
	}

	logger.Logger().Info("UpdateConfig 出参", zap.Any("data", config))
	utils.JSON(c, utils.Success(config))
}

// DELETE /api/system/config/:id
func (h *ConfigHandler) DeleteConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	logger.Logger().Info("DeleteConfig 入参", zap.Int64("id", id))

	if err := h.service.Delete(id); err != nil {
		logger.Logger().Error("DeleteConfig 失败", zap.Error(err))
		utils.JSON(c, utils.Error(err.Error(), http.StatusInternalServerError))
		return
	}

	logger.Logger().Info("DeleteConfig 出参", zap.String("msg", "删除成功"))
	utils.JSON(c, utils.Success("删除成功"))
}

func init() {
	// 自动注册路由模块（通过 init 自动调用）
	router.RegisterRouteModule(&ConfigHandler{})
}

func (h *ConfigHandler) Register(rg *gin.RouterGroup, db *gorm.DB) {
	h.service = service.NewConfigService(repository.NewConfigRepository(db))
	api := rg.Group("/system/config")
	{
		api.GET("/list", h.GetConfigList)
		api.GET("/:id", h.GetConfigById)
		api.POST("", h.AddConfig)
		api.PUT("", h.UpdateConfig)
		api.DELETE("/:id", h.DeleteConfig)
	}

}
