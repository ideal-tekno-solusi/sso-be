package handler

import (
	"app/api/sso/operation"
	"app/internal/sso/entity"
	"app/internal/sso/repository"
	"app/utils"
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (r *RestService) Authorize(ctx echo.Context, params *operation.AuthorizeRequest) error {
	context := context.Background()

	repo := repository.InitRepo(r.dbr, r.dbw)
	authorizeService := repository.AuthorizeRepository(repo)

	//? used for validation, continue if exist and valid
	if params.ResponseType == "refresh" {
		tokenValid, err := utils.ValidateJwt(params.State)
		if err != nil {
			errorMessage := fmt.Sprintf("failed to validate token with error: %v", err)
			logrus.Error(errorMessage)

			utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

			return nil
		}
		if !tokenValid {
			errorMessage := "token is not valid, please try login again"
			logrus.Warn(errorMessage)

			utils.SendProblemDetailJson(ctx, http.StatusUnauthorized, errorMessage, ctx.Path(), uuid.NewString())

			return nil
		}

		refreshToken, err := authorizeService.GetRefreshToken(context, params.State)
		if err != nil {
			errorMessage := fmt.Sprintf("failed to get refresh token with error: %v", err)
			logrus.Error(errorMessage)

			utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

			return nil
		}
		if refreshToken == nil {
			errorMessage := "refresh token not found, please try login again"
			logrus.Warn(errorMessage)

			utils.SendProblemDetailJson(ctx, http.StatusUnauthorized, errorMessage, ctx.Path(), uuid.NewString())

			return nil
		}

		err = authorizeService.CreateSession(context, params.State, refreshToken.UserID.String, params.ClientId, params.CodeChallenge, params.CodeChallengeMethod, params.Scopes, params.RedirectUrl)
		if err != nil {
			errorMessage := fmt.Sprintf("failed to create new session with error: %v", err)
			logrus.Error(errorMessage)

			utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

			return nil
		}

		err = authorizeService.DeleteRefreshToken(context, params.State)
		if err != nil {
			errorMessage := fmt.Sprintf("failed to delete refresh token with error: %v", err)
			logrus.Error(errorMessage)

			utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

			return nil
		}
	}

	//? used for validation, continue if exist
	if params.ResponseType == "code" {
		session, err := authorizeService.GetSession(context, params.State)
		if err != nil {
			errorMessage := fmt.Sprintf("failed to get session from db with error: %v", err)
			logrus.Error(errorMessage)

			utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

			return nil
		}
		if session == nil {
			errorMessage := "session not found in db, please try again later"
			logrus.Warn(errorMessage)

			utils.SendProblemDetailJson(ctx, http.StatusForbidden, errorMessage, ctx.Path(), uuid.NewString())

			return nil
		}
	}

	authorizeCode, err := utils.GenerateRandomString(64)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to generate auth code with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	err = authorizeService.CreateAuthToken(context, *authorizeCode, params.State)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to insert auth token with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	res := entity.AuthorizeResponse{
		AuthorizeCode: *authorizeCode,
	}

	return ctx.JSON(http.StatusOK, utils.GenerateResponseJson(nil, true, res))
}
