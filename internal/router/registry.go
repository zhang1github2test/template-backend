package router

import "github.com/gin-gonic/gin"

type RouteRegistrar interface {
	Register(rg *gin.RouterGroup)
}

var registrars []RouteRegistrar

func RegisterRouteModule(module RouteRegistrar) {
	registrars = append(registrars, module)
}

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	for _, m := range registrars {
		m.Register(api)
	}
}
