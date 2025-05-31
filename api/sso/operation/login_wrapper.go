package operation

import (
	"app/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Username            string `json:"username" validate:"required"`
	Password            string `json:"password" validate:"required"`
	RedirectUrl         string `json:"redirectUrl" validate:"required"`
	ClientId            string `json:"clientId" validate:"required"`
	ResponseType        string `json:"responseType" validate:"required,oneofci=code refresh"`
	Scopes              string `json:"scopes"`
	State               string `json:"state" validate:"required"`
	CodeChallenge       string `json:"codeChallenge" validate:"required"`
	CodeChallengeMethod string `json:"codeChallengeMethod" validate:"required,eq=S256"`
}

func LoginWrapper(handler func(e echo.Context, params *LoginRequest) error) echo.HandlerFunc {
	return func(e echo.Context) error {
		params := LoginRequest{}

		err := e.Bind(&params)
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
