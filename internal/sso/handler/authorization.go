package handler

import (
	"app/api/sso/operation"
	"app/internal/sso/repository"
	"app/utils"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (r *RestService) Authorization(ctx echo.Context, params *operation.AuthorizationRequest) error {
	//TODO: lanjutin bikin logic nya untuk simpan code challenge ke db dan cek cookie authorization apakah exist
	//TODO: sementara session id generate dari uuid, harap gunakan teknik lain
	context := context.Background()

	repo := repository.InitRepo(r.dbr, r.dbw)
	authorizationService := repository.AuthorizationRepository(repo)

	authorization, _ := ctx.Cookie("SSO-AUTHORIZATION")
	if authorization != nil {
		token, _ := url.QueryUnescape(authorization.Value)

		data, err := authorizationService.GetAuthorization(context, token)
		if err != nil {
			errorMessage := fmt.Sprintf("failed to get authorization token with error: %v", err)
			logrus.Warn(errorMessage)

			utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

			return nil
		}
		if data == nil {
			errorMessage := "current authorization token not found, please relog and try again"
			logrus.Info(errorMessage)

			//? delete authorization cookie
			utils.SetCookie(ctx, http.Cookie{
				Name:     "SSO-AUTHORIZATION",
				Path:     "/",
				Domain:   "localhost",
				MaxAge:   -1,
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
			})

			utils.SendProblemDetailJson(ctx, http.StatusUnauthorized, errorMessage, ctx.Path(), uuid.NewString())

			return nil
		}

		return ctx.Redirect(http.StatusPermanentRedirect, params.RedirectUrl)
	}

	age := viper.GetInt("config.session.age")
	domain := viper.GetString("config.session.domain")
	path := viper.GetString("config.session.path")

	sessionId := uuid.NewString()

	utils.SetCookie(ctx, http.Cookie{
		Name:     "Session-Id",
		Value:    sessionId,
		MaxAge:   age,
		Domain:   domain,
		Path:     path,
		Secure:   false,
		HttpOnly: true,
	})

	err := authorizationService.CreateSession(context, sessionId, params.ClientId, params.CodeChallenge, params.CodeChallengeMethod)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to create new session with error: %v", err)
		logrus.Warn(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//TODO: redirect to sso fe login page

	//? set back csrf from sender to header
	ctx.Response().Header().Set("X-CSRF-Token", params.State)

	return ctx.NoContent(200)
}
