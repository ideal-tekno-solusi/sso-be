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
)

func (r *RestService) Authorization(ctx echo.Context, params *operation.AuthorizationRequest) error {
	//TODO: lanjutin bikin logic nya untuk simpan code challenge ke db dan cek cookie authorization apakah exist
	//TODO: sementara session id generate dari uuid, harap gunakan teknik lain
	context := context.Background()

	repo := repository.InitRepo(r.dbr, r.dbw)
	authorizationService := repository.AuthorizationRepository(repo)

	sessionId := uuid.NewString()

	cookie := new(http.Cookie)
	cookie.Name = "Session-Id"
	cookie.Value = sessionId
	cookie.MaxAge = 3600
	cookie.Domain = "localhost"
	cookie.Path = "/"
	cookie.Secure = false
	cookie.HttpOnly = true
	ctx.SetCookie(cookie)

	err := authorizationService.CreateSession(context, sessionId, params.ClientId, params.CodeChallenge, params.CodeChallengeMethod)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to create new session with error: %v", err)
		logrus.Warn(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())
	}

	return ctx.NoContent(200)
}
