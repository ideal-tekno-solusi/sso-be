package handler

import (
	"app/api/sso/operation"
	"app/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (r *RestService) User(ctx echo.Context, params *operation.UserRequest) error {
	serviceName := "GET User"

	if params.Token == "" {
		//? get session, will create new if not found
		sess, err := session.Get("session", ctx)
		if err != nil {
			errorMessage := fmt.Sprintf("failed to get session with error: %v", err)
			utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

			utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

			return nil
		}
		if sess.IsNew {
			errorMessage := "failed to process user because token not found in payload and session"
			utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

			utils.SendProblemDetailJson(ctx, http.StatusUnauthorized, errorMessage, ctx.Path(), uuid.NewString())

			return nil
		}

		accessToken := sess.Values["access_token"]
		if accessToken == nil {
			errorMessage := "access token not found"
			utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

			utils.SendProblemDetailJson(ctx, http.StatusUnauthorized, errorMessage, ctx.Path(), uuid.NewString())

			return nil
		}

		params.Token = fmt.Sprintf("Bearer %v", accessToken.(string))
	}

	token := strings.Split(params.Token, " ")

	plainText, err := utils.DecryptJwt(token[1])
	if err != nil {
		errorMessage := fmt.Sprintf("failed to decrypt jwt with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusUnauthorized, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	return ctx.JSON(http.StatusOK, utils.GenerateResponseJson(nil, true, plainText))
}
