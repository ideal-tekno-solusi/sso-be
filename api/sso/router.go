package sso

import (
	"app/api/sso/operation"

	"github.com/labstack/echo/v4"
)

func Router(r *echo.Echo, s Service) {
	v1 := r.Group("/v1")
	v1.GET("/api/authorization", operation.AuthorizationWrapper(s.Authorization))
}
