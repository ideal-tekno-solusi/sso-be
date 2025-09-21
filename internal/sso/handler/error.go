package handler

import (
	"app/api/sso/operation"
	"app/internal/sso/entity"
	"app/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *RestService) Error(ctx echo.Context, params *operation.ErrorRequest) error {
	// context := context.Background()
	// serviceName := "GET Error"

	res := entity.ErrorResponse{
		Message: utils.ErrorCodeString(params.Code),
	}

	return ctx.JSON(http.StatusOK, utils.GenerateResponseJson(nil, true, res))
}
