// internal/handler/menu_handler.go
package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"template-backend/internal/model"
	"template-backend/internal/repository"
	"template-backend/internal/router"
	"template-backend/internal/service"
	"template-backend/pkg/logger"
	"template-backend/pkg/utils"
)

type MenuHandler struct {
	service *service.MenuService
}

func NewMenuHandler(service *service.MenuService) *MenuHandler {
	return &MenuHandler{service: service}
}

func (h *MenuHandler) GetMenuTree(c *gin.Context) {
	menu := model.Menu{}
	menuType := c.Query("type")
	if menuType != "" {
		menuTypeInt, err := strconv.Atoi(menuType)
		if err != nil {
			utils.JSON(c, utils.Error("invalid menu type", http.StatusBadRequest))
			return
		}
		menu.Type = menuTypeInt
	}

	name := c.Query("name")
	if name != "" {
		menu.Name = name
	}
	visible := c.Query("visible")
	if visible != "" {
		parseBool, err := strconv.ParseBool(visible)
		if err != nil {
			utils.JSON(c, utils.Error("invalid visible", http.StatusBadRequest))
			return
		}
		menu.Visible = &parseBool
	}

	menus, err := h.service.GetMenuTree(menu)
	if err != nil {
		utils.JSON(c, utils.Error("get menu tree failed: "+err.Error(), http.StatusInternalServerError))
		return
	}
	for i := range menus {
		menus[i].UnMarshalMeta()
	}
	utils.JSON(c, utils.Success(menus))
}

func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var menu model.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.Logger().Info("menu", zap.String("name", menu.Name))
	menu.MarshalMeta()
	logger.Logger().Info("menu", zap.String("name", menu.Name))
	if err := h.service.CreateMenu(&menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSON(c, utils.Success(menu))
}

func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var menu model.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	menu.ID = uint(id)
	menu.MarshalMeta()

	if err := h.service.UpdateMenu(&menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSON(c, utils.Success(menu))
}

func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteMenu(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSON(c, utils.Success(gin.H{"message": "deleted"}))
}

func (h *MenuHandler) Register(rg *gin.RouterGroup, db *gorm.DB) {
	h.service = service.NewMenuService(repository.NewMenuRepository(db))
	menu := rg.Group("/menu")
	{
		menu.GET("/tree", h.GetMenuTree)
		menu.POST("", h.CreateMenu)
		menu.PUT("/:id", h.UpdateMenu)
		menu.DELETE("/:id", h.DeleteMenu)
	}
}

func init() {
	// 自动注册路由模块（通过 init 自动调用）
	router.RegisterRouteModule(&MenuHandler{})
}
