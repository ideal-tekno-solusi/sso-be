package sso

import (
	"app/api/sso/operation"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Router(r *echo.Echo, s Service) {
	auth := r.Group("/auth")
	auth.GET("/api/authorize", operation.AuthorizeWrapper(s.Authorize))
	auth.POST("/api/login", operation.LoginWrapper(s.Login))
	auth.POST("/api/token", operation.TokenWrapper(s.Token))

	test := r.Group("/test")
	test.GET("/redirect", func(c echo.Context) error {
		c.Redirect(http.StatusPermanentRedirect, "https://inventory.idtecsi.my.id/test/redirect")
		return nil
	})
}
