package api

import (
	"app/bootstrap"
	sso "app/internal/sso/handler"

	"github.com/labstack/echo/v4"
)

func RegisterApi(r *echo.Echo, cfg *bootstrap.Container) {
	sso.RestRegister(r, cfg)
}
