package operation

import (
	"app/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ErrorRequest struct {
	Code string `query:"code" validate:"required"`
}

func ErrorWrapper(handler func(e echo.Context, params *ErrorRequest) error) echo.HandlerFunc {
	return func(e echo.Context) error {
		params := ErrorRequest{}

		err := (&echo.DefaultBinder{}).BindQueryParams(e, &params)
		if err != nil {
			utils.SendProblemDetailJson(e, http.StatusInternalServerError, err.Error(), e.Path(), uuid.NewString())

			return nil
		}

		err = e.Validate(params)
		if err != nil {
			utils.SendProblemDetailJsonValidate(e, http.StatusBadRequest, "validation error", e.Path(), uuid.NewString(), err.(validator.ValidationErrors))

			return nil
		}

		return handler(e, &params)
	}
}
