package user

import (
	"bruno_authentication/pkg/infrastructure"
	"bruno_authentication/pkg/middlewares"
)

type Route struct {
	router               *infrastructure.Router
	controller           *Controller
	permissionMiddleware middlewares.PermissionMiddleware
}

func NewRoute(router *infrastructure.Router, controller *Controller, permissionMiddleware middlewares.PermissionMiddleware) *Route {
	route := Route{router: router, controller: controller, permissionMiddleware: permissionMiddleware}
	return &route
}

func RegisterRoute(r *Route) {
	basicAuthRoutes := r.router.Group("/auth", r.permissionMiddleware.BasicAuthPermission())
	basicAuthRoutes.POST("/login", r.controller.Login)
}
