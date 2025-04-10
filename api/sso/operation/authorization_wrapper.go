package operation

import (
	"app/utils"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthorizationRequest struct {
	RedirectUrl         string `query:"redirect_url"`
	ClientId            string `query:"client_id"`
	ResponseType        string `query:"response_type"`
	Scopes              string `query:"scopes"`
	State               string `query:"state"`
	CodeChallenge       string `query:"code_challenge"`
	CodeChallengeMethod string `query:"code_challenge_method"`
}

func AuthorizationWrapper(handler func(e echo.Context, params *AuthorizationRequest) error) echo.HandlerFunc {
	return func(e echo.Context) error {
		params := AuthorizationRequest{}

		err := (&echo.DefaultBinder{}).BindQueryParams(e, &params)
		if err != nil {
			utils.SendProblemDetailJson(e, http.StatusInternalServerError, err.Error(), e.Path(), uuid.NewString())

			return nil
		}

		err = validateAuthorizationRequest(params)
		if err != nil {
			utils.SendProblemDetailJson(e, http.StatusBadRequest, err.Error(), e.Path(), uuid.NewString())

			return nil
		}

		csrfToken, _ := e.Cookie(fmt.Sprintf("%v-XSRF-TOKEN", strings.ToUpper(params.ClientId)))
		if csrfToken == nil {
			err := errors.New("xsrf token is empty")
			utils.SendProblemDetailJson(e, http.StatusBadRequest, err.Error(), e.Path(), uuid.NewString())

			return nil
		}

		token, _ := url.QueryUnescape(csrfToken.Value)
		if params.State != token {
			err := errors.New("state and xsrf token is not match")
			utils.SendProblemDetailJson(e, http.StatusBadRequest, err.Error(), e.Path(), uuid.NewString())

			return nil
		}

		return handler(e, &params)
	}
}

func validateAuthorizationRequest(params AuthorizationRequest) error {
	if len(params.ClientId) == 0 || strings.TrimSpace(params.ClientId) == "" {
		return errors.New("client id can't be empty")
	}

	if len(params.RedirectUrl) == 0 || strings.TrimSpace(params.RedirectUrl) == "" {
		return errors.New("redirect url can't be empty")
	}

	if len(params.ResponseType) == 0 || strings.TrimSpace(params.ResponseType) == "" {
		return errors.New("response type can't be empty")
	}

	if strings.ToLower(params.ResponseType) != "code" {
		return errors.New("response type need to be set to code")
	}

	if len(params.State) == 0 || strings.TrimSpace(params.State) == "" {
		return errors.New("state can't be empty")
	}

	if len(params.CodeChallenge) == 0 || strings.TrimSpace(params.CodeChallenge) == "" {
		return errors.New("code challenge can't be empty")
	}

	if len(params.CodeChallengeMethod) == 0 || strings.TrimSpace(params.CodeChallengeMethod) == "" {
		return errors.New("code challenge method can't be empty")
	}

	if params.CodeChallengeMethod != "S256" {
		return errors.New("code challenge method only accept S256")
	}

	return nil
}
