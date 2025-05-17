package handler

import (
	"app/api/sso/operation"
	"app/internal/sso/repository"
	"app/utils"
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (r *RestService) Authorization(ctx echo.Context, params *operation.AuthorizationRequest) error {
	context := context.Background()

	repo := repository.InitRepo(r.dbr, r.dbw)
	authorizationService := repository.AuthorizationRepository(repo)

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

	err := authorizationService.CreateSession(context, sessionId, params.ClientId, params.CodeChallenge, params.CodeChallengeMethod, params.Scopes)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to create new session with error: %v", err)
		logrus.Warn(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//TODO: redirect to sso fe login page

	return ctx.NoContent(200)
}
