package sso

import (
	"app/api/sso/operation"

	"github.com/labstack/echo/v4"
)

type Service interface {
	Authorize(ctx echo.Context, params *operation.AuthorizeRequest) error
	Login(ctx echo.Context, params *operation.LoginRequest) error
	Token(ctx echo.Context, params *operation.TokenRequest) error
	User(ctx echo.Context, params *operation.UserRequest) error
	Error(ctx echo.Context, params *operation.ErrorRequest) error
}
