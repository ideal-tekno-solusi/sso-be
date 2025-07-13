package operation

import (
	"app/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// type TokenRequest struct {
// 	Code         string `json:"code" validate:"required"`
// 	CodeVerifier string `json:"codeVerifier" validate:"required"`
// }

type TokenRequest struct {
	GrantType    string `json:"grant_type" validate:"required"`
	Code         string `json:"code"`
	RefreshToken string `json:"refresh_code"`
	RedirectUri  string `json:"redirect_uri"`
	CodeVerifier string `json:"code_verifier"`
	ClientId     string `json:"client_id"`
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
