package operation

import (
	"app/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthorizeRequest struct {
	RedirectUrl         string `json:"redirect_url" query:"redirect_url" validate:"required"`
	ClientId            string `json:"client_id" query:"client_id" validate:"required"`
	ResponseType        string `json:"response_type" query:"response_type" validate:"required,oneofci=code refresh"`
	Scopes              string `json:"scopes" query:"scopes"`
	State               string `json:"state" query:"state" validate:"required"`
	CodeChallenge       string `json:"code_challenge" query:"code_challenge" validate:"required"`
	CodeChallengeMethod string `json:"code_challenge_method" query:"code_challenge_method" validate:"required,eq=S256"`
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
