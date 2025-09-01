package sso

import (
	"app/api/sso/operation"

	"github.com/labstack/echo/v4"
)

func Router(r *echo.Echo, s Service) {
	api := r.Group("/api")
	api.GET("/v1/authorize", operation.AuthorizeWrapper(s.Authorize))
	api.POST("/v1/login", operation.LoginWrapper(s.Login))
	api.POST("/v1/token", operation.TokenWrapper(s.Token))
	api.GET("/v1/user", operation.UserWrapper(s.User))
}
