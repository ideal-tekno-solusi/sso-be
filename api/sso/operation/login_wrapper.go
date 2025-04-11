package operation

import (
	"app/utils"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginWrapper(handler func(e echo.Context, params *LoginRequest) error) echo.HandlerFunc {
	return func(e echo.Context) error {
		params := LoginRequest{}

		err := e.Bind(&params)
		if err != nil {
			utils.SendProblemDetailJson(e, http.StatusInternalServerError, err.Error(), e.Path(), uuid.NewString())

			return nil
		}

		err = (&echo.DefaultBinder{}).BindHeaders(e, &params)
		if err != nil {
			utils.SendProblemDetailJson(e, http.StatusInternalServerError, err.Error(), e.Path(), uuid.NewString())

			return nil
		}

		err = validateLoginRequest(params)
		if err != nil {
			utils.SendProblemDetailJson(e, http.StatusBadRequest, err.Error(), e.Path(), uuid.NewString())

			return nil
		}

		return handler(e, &params)
	}
}

func validateLoginRequest(params LoginRequest) error {
	if len(params.Username) == 0 || strings.TrimSpace(params.Username) == "" {
		return errors.New("username can't be empty")
	}

	if len(params.Password) == 0 || strings.TrimSpace(params.Password) == "" {
		return errors.New("password can't be empty")
	}

	return nil
}
