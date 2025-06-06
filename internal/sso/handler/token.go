package handler

import (
	"app/api/sso/operation"
	"app/internal/sso/entity"
	"app/internal/sso/repository"
	"app/utils"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (r *RestService) Token(ctx echo.Context, params *operation.TokenRequest) error {
	context := context.Background()

	repo := repository.InitRepo(r.dbr, r.dbw)
	tokenService := repository.TokenRepository(repo)

	//? generate code challenge from code verifier req
	hash := sha256.New()
	hash.Write([]byte(params.CodeVerifier))
	codeChallenge := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	//? get token by code challenge and check validity of code verifier at once
	token, err := tokenService.GetToken(context, codeChallenge)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to get token with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? check if auth code is match
	if params.Code != token.ID {
		errorMessage := "token not match, please login again later"
		logrus.Warn(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusUnauthorized, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? generate jwt token
	jwtBody := map[string]string{
		"username":     token.Username,
		"name":         token.Name,
		"redirect_url": token.RedirectUrl.String,
	}

	tokenExpTime := viper.GetInt("secret.expToken")
	refreshExpTime := viper.GetInt("secret.refreshToken")

	accessToken, err := utils.GenerateAuthToken(jwtBody, tokenExpTime)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to generate access token with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	refreshToken, err := utils.GenerateAuthToken(jwtBody, refreshExpTime)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to generate access token with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	err = tokenService.CreateRefreshToken(context, *refreshToken, token.Username)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to insert refresh token with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	authToken := entity.Token{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
		ExpiresIn:    tokenExpTime,
		Scope:        token.Scopes.String,
		TokenType:    "Bearer",
	}

	//? delete auth token and session and remove session id
	err = tokenService.DeleteAuthToken(context, token.SessionID.String)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to delete auth token with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	err = tokenService.DeleteSession(context, token.SessionID.String)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to delete session with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	return ctx.JSON(http.StatusOK, authToken)
}
