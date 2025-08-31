package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RouteRegistrar interface {
	Register(rg *gin.RouterGroup, db *gorm.DB)
}

var registrars []RouteRegistrar

func RegisterRouteModule(module RouteRegistrar) {
	registrars = append(registrars, module)
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api")
	for _, m := range registrars {
		m.Register(api, db)
	}
}
