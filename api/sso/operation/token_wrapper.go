package operation

import (
	"app/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TokenRequest struct {
	GrantType    string `json:"grant_type" validate:"required,oneofci=authorization_code refresh"`
	Code         string `json:"code"`
	RedirectUrl  string `json:"redirect_url"`
	CodeVerifier string `json:"code_verifier"`
	ClientId     string `json:"client_id" validate:"required"`
	ClientSecret string `json:"client_secret"`
}

func TokenWrapper(handler func(e echo.Context, params *TokenRequest) error) echo.HandlerFunc {
	return func(e echo.Context) error {
		params := TokenRequest{}

		err := (&echo.DefaultBinder{}).BindBody(e, &params)
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
