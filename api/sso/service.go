package sso

import (
	"app/api/sso/operation"

	"github.com/labstack/echo/v4"
)

type Service interface {
	Authorization(ctx echo.Context, params *operation.AuthorizationRequest) error
}
