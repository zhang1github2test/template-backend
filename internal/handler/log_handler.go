package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"template-backend/internal/repository"
	"template-backend/internal/router"
	"template-backend/internal/service"
	"template-backend/pkg/logger"
	"template-backend/pkg/utils"
	"time"
)

type LogHandler interface {
	GetLogList(c *gin.Context)
	GetLogByID(c *gin.Context)
	DeleteLog(c *gin.Context)
	DeleteLogs(c *gin.Context)
	CleanLogs(c *gin.Context)
	ExportLogs(c *gin.Context)
}

type logHandler struct {
	service service.LogService
}

func NewLogHandler(service service.LogService) LogHandler {
	return &logHandler{service: service}
}

// GetLogList 获取日志列表
func (h *logHandler) GetLogList(c *gin.Context) {
	logger.Logger().Info("Handling GetLogList request")

	// 解析查询参数
	pageNum := 1
	if pageNumStr := c.Query("pageNum"); pageNumStr != "" {
		if n, err := strconv.Atoi(pageNumStr); err == nil && n > 0 {
			pageNum = n
		}
	}

	pageSize := 10
	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		if n, err := strconv.Atoi(pageSizeStr); err == nil && n > 0 {
			pageSize = n
		}
	}

	// 构建查询条件
	conditions := make(map[string]interface{})
	if method := c.Query("method"); method != "" {
		conditions["method"] = method
	}
	if path := c.Query("path"); path != "" {
		conditions["path"] = path
	}
	if statusStr := c.Query("status"); statusStr != "" {
		if status, err := strconv.Atoi(statusStr); err == nil {
			conditions["status"] = status
		}
	}
	if ip := c.Query("ip"); ip != "" {
		conditions["ip"] = ip
	}
	if handler := c.Query("handler"); handler != "" {
		conditions["handler"] = handler
	}
	if timestamps := c.QueryArray("timestamp[]"); len(timestamps) == 2 {
		// 设置东八区时区
		loc, _ := time.LoadLocation("Asia/Shanghai")
		startTime, err1 := time.ParseInLocation(time.DateTime, timestamps[0], loc)
		endTime, err2 := time.ParseInLocation(time.DateTime, timestamps[1], loc)
		if err1 == nil && err2 == nil {
			conditions["timestamp"] = []time.Time{startTime, endTime}
		}
	}

	logs, total, err := h.service.GetLogList(pageNum, pageSize, conditions)
	if err != nil {
		logger.Logger().Error("Failed to get log list", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取日志列表失败",
			"data": nil,
		})
		return
	}
	utils.JSON(c, utils.Success(gin.H{
		"rows":  logs,
		"total": total,
		"page":  pageNum,
		"size":  pageSize,
	}))
}

// GetLogByID 根据ID获取日志详情
func (h *logHandler) GetLogByID(c *gin.Context) {
	logger.Logger().Info("Handling GetLogByID request")

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Logger().Error("Invalid log ID", zap.String("id", idStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的日志ID",
			"data": nil,
		})
		return
	}

	log, err := h.service.GetLogByID(uint(id))
	if err != nil {
		logger.Logger().Error("Failed to get log by ID", zap.Uint("id", uint(id)), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "日志不存在",
			"data": nil,
		})
		return
	}

	utils.JSON(c, utils.Success(log))
}

// DeleteLog 删除日志
func (h *logHandler) DeleteLog(c *gin.Context) {
	logger.Logger().Info("Handling DeleteLog request")

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Logger().Error("Invalid log ID", zap.String("id", idStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的日志ID",
			"data": nil,
		})
		return
	}

	err = h.service.DeleteLog(uint(id))
	if err != nil {
		logger.Logger().Error("Failed to delete log", zap.Uint("id", uint(id)), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除日志失败",
			"data": nil,
		})
		return
	}
	utils.JSON(c, utils.Success(nil))
}

// DeleteLogs 批量删除日志
func (h *logHandler) DeleteLogs(c *gin.Context) {
	logger.Logger().Info("Handling DeleteLogs request")

	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger().Error("Failed to bind delete logs request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
			"data": nil,
		})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请选择要删除的日志",
			"data": nil,
		})
		return
	}

	err := h.service.DeleteLogs(req.IDs)
	if err != nil {
		logger.Logger().Error("Failed to batch delete logs", zap.Any("ids", req.IDs), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "批量删除日志失败",
			"data": nil,
		})
		return
	}

	utils.JSON(c, utils.Success(nil))
}

// CleanLogs 清空日志
func (h *logHandler) CleanLogs(c *gin.Context) {
	logger.Logger().Info("Handling CleanLogs request")

	err := h.service.CleanLogs()
	if err != nil {
		logger.Logger().Error("Failed to clean logs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "清空日志失败",
			"data": nil,
		})
		return
	}

	utils.JSON(c, utils.Success(nil))
}

// ExportLogs 导出日志
func (h *logHandler) ExportLogs(c *gin.Context) {
	logger.Logger().Info("Handling ExportLogs request")

	// 这里实现日志导出逻辑，可以导出为CSV或Excel格式
	// 为简化示例，这里返回一个简单的文本文件

	content := "日志导出功能待实现"

	c.Header("Content-Type", "text/plain")
	c.Header("Content-Disposition", "attachment; filename=logs.txt")
	c.Data(http.StatusOK, "text/plain", []byte(content))
}

func init() {
	// 自动注册路由模块（通过 init 自动调用）
	router.RegisterRouteModule(&logHandler{})
}

func (h *logHandler) Register(rg *gin.RouterGroup, db *gorm.DB) {
	h.service = service.NewLogService(repository.NewLogRepository(db))
	// 日志相关路由
	logGroup := rg.Group("/system/log")
	{
		logGroup.GET("/list", h.GetLogList)
		logGroup.GET("/:id", h.GetLogByID)
		logGroup.DELETE("/:id", h.DeleteLog)
		logGroup.DELETE("", h.DeleteLogs)
		logGroup.DELETE("/clean", h.CleanLogs)
		logGroup.POST("/export", h.ExportLogs)
	}

}
