package handler

import (
	"app/api/sso/operation"
	"app/internal/sso/entity"
	"app/internal/sso/repository"
	"app/utils"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (r *RestService) Token(ctx echo.Context, params *operation.TokenRequest) error {
	context := context.Background()

	repo := repository.InitRepo(r.dbr, r.dbw)
	tokenService := repository.TokenRepository(repo)

	sessionId, err := ctx.Cookie("Session-Id")
	if err != nil {
		errorMessage := "session id not found, please login again later"
		logrus.Warn(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusForbidden, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	token, err := tokenService.GetToken(context, sessionId.String())
	if err != nil {
		errorMessage := fmt.Sprintf("failed to get token with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? generate code challenge from code verifier req
	hash := sha256.New()
	hash.Write([]byte(params.CodeVerifier))
	codeChallenge := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	//? check is generated code challenge same asa saved code challenge
	if codeChallenge != token.CodeChallenge {
		errorMessage := "code challenge is not same, please login again later"
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusUnauthorized, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? check if auth code is match
	if params.Code != token.ID {
		errorMessage := "token not match, please login again later"
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusUnauthorized, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? generate jwt token
	jwtBody := entity.Jwt{
		Name: token.Name,
	}

	jwtBodyString, _ := json.Marshal(jwtBody)
	tokenExpTime := viper.GetInt("secret.expToken")
	refreshExpTime := viper.GetInt("secret.refreshToken")

	accessToken, err := utils.GenerateAuthToken(string(jwtBodyString), token.Username, tokenExpTime)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to generate access token with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	refreshToken, err := utils.GenerateAuthToken(string(jwtBodyString), token.Username, refreshExpTime)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to generate access token with error: %v", err)
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
	err = tokenService.DeleteAuthToken(context, sessionId.String())
	if err != nil {
		errorMessage := fmt.Sprintf("failed to delete auth token with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	err = tokenService.DeleteSession(context, sessionId.String())
	if err != nil {
		errorMessage := fmt.Sprintf("failed to delete session with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	utils.SetCookie(ctx, http.Cookie{
		Name:     "Session-Id",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		Secure:   false,
		HttpOnly: true,
	})

	return ctx.JSON(http.StatusOK, authToken)
}
