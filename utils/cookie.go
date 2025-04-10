package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetCookie(ctx echo.Context, cookie http.Cookie) {
	ctx.SetCookie(&cookie)
}
