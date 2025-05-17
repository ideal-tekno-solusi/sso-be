package sso

import (
	"app/api/sso/operation"

	"github.com/labstack/echo/v4"
)

func Router(r *echo.Echo, s Service) {
	auth := r.Group("/auth")
	auth.GET("/api/authorization", operation.AuthorizationWrapper(s.Authorization))
	auth.POST("/api/login", operation.LoginWrapper(s.Login))
	auth.POST("/api/token", operation.TokenWrapper(s.Token))
}
