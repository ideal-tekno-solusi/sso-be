package operation

import (
	"app/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthorizeRequest struct {
	RedirectUrl         string `query:"redirectUrl" validate:"required"`
	ClientId            string `query:"clientId" validate:"required"`
	ResponseType        string `query:"responseType" validate:"required,oneofci=code refresh"`
	Scopes              string `query:"scopes"`
	State               string `query:"state" validate:"required"`
	CodeChallenge       string `query:"codeChallenge" validate:"required"`
	CodeChallengeMethod string `query:"codeChallengeMethod" validate:"required,eq=S256"`
}

func AuthorizeWrapper(handler func(e echo.Context, params *AuthorizeRequest) error) echo.HandlerFunc {
	return func(e echo.Context) error {
		params := AuthorizeRequest{}

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
