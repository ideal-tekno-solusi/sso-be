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
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (r *RestService) Token(ctx echo.Context, params *operation.TokenRequest) error {
	context := context.Background()
	serviceName := "POST Token"

	repo := repository.InitRepo(r.dbr, r.dbw)
	tokenService := repository.TokenRepository(repo)

	//? get session, will create new if not found
	sess, err := session.Get("session", ctx)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to get session with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}
	if sess.IsNew {
		errorMessage := "wrong request order, please login by click login button in registered website first"
		utils.WarningLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusUnauthorized, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? check if user already login
	guid := sess.Values["guid"]
	if guid == nil {
		errorMessage := "current user not login yet, please try to login again"
		utils.WarningLog(errorMessage, ctx.Path(), serviceName)
		utils.DeleteSession(sess, ctx.Request(), ctx.Response())

		utils.SendProblemDetailJson(ctx, http.StatusUnauthorized, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? check if client exist first
	client, err := tokenService.GetClient(context, params.ClientId)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to get client from db with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}
	if client == nil {
		errorMessage := "client not found"
		utils.WarningLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusForbidden, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}
	if client.Type == "SS" {
		if client.Secret.String != params.ClientSecret {
			errorMessage := "client id or secret not valid"
			utils.WarningLog(errorMessage, ctx.Path(), serviceName)

			utils.SendProblemDetailJson(ctx, http.StatusForbidden, errorMessage, ctx.Path(), uuid.NewString())

			return nil
		}
	}

	//? run this flow if grant type is refresh and client is SPA and refresh token not provided in req
	if params.GrantType == "refresh" {
		if params.Code == "" && client.Type == "SPA" {
			refreshToken := sess.Values["refresh_token"]
			if refreshToken == nil {
				errorMessage := "refresh token not found, please login again"
				utils.WarningLog(errorMessage, ctx.Path(), serviceName)

				utils.SendProblemDetailJson(ctx, http.StatusUnauthorized, errorMessage, ctx.Path(), uuid.NewString())

				return nil
			}

			params.Code = refreshToken.(string)
		}
	}

	if params.GrantType == "authorization_code" {
		//? generate code challenge from code verifier req
		codeChallengeSource := sess.Values["code_Challenge"]
		codeChallengeMethodSource := sess.Values["code_challenge_method"]
		if codeChallengeSource == nil || codeChallengeMethodSource == nil {
			errorMessage := "code challenge not found, please try to login again"
			utils.WarningLog(errorMessage, ctx.Path(), serviceName)
			utils.DeleteSession(sess, ctx.Request(), ctx.Response())

			utils.SendProblemDetailJson(ctx, http.StatusUnauthorized, errorMessage, ctx.Path(), uuid.NewString())

			return nil
		}

		var codeChallenge string = params.CodeVerifier

		if codeChallengeMethodSource == "S256" {
			hash := sha256.New()
			hash.Write([]byte(params.CodeVerifier))
			codeChallenge = base64.StdEncoding.EncodeToString(hash.Sum(nil))
		}

		//? check if code verifier equal to code challenge
		if codeChallengeSource.(string) != codeChallenge {
			errorMessage := "code verifier not match"
			utils.WarningLog(errorMessage, ctx.Path(), serviceName)

			utils.SendProblemDetailJson(ctx, http.StatusForbidden, errorMessage, ctx.Path(), uuid.NewString())

			return nil
		}
	}

	//? validate code and user id of current session
	auth, err := tokenService.GetAuth(context, params.Code)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to get auth code with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}
	if auth == nil {
		errorMessage := "auth code not found, please try login again"
		utils.WarningLog(errorMessage, ctx.Path(), serviceName)
		utils.DeleteSession(sess, ctx.Request(), ctx.Response())

		utils.SendProblemDetailJson(ctx, http.StatusUnauthorized, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	session, err := tokenService.GetSession(context, guid.(string))
	if err != nil {
		errorMessage := fmt.Sprintf("failed to get session with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}
	if session == nil {
		errorMessage := "session not found, please try login again"
		utils.WarningLog(errorMessage, ctx.Path(), serviceName)
		utils.DeleteSession(sess, ctx.Request(), ctx.Response())

		utils.SendProblemDetailJson(ctx, http.StatusUnauthorized, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}
	if session.UserID.String != auth.UserID.String {
		errorMessage := "user of current auth is not same with current session, please try login again"
		utils.WarningLog(errorMessage, ctx.Path(), serviceName)
		utils.DeleteSession(sess, ctx.Request(), ctx.Response())

		utils.SendProblemDetailJson(ctx, http.StatusUnauthorized, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? generate jwt token
	jwtBody := map[string]string{
		"username": session.UserID.String,
	}

	tokenExpTime := client.TokenLivetime.Int64

	accessToken, err := utils.GenerateAuthToken(jwtBody, int(tokenExpTime))
	if err != nil {
		errorMessage := fmt.Sprintf("failed to generate access token with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	refreshToken, err := utils.GenerateRandomString(64)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to generate refresh token with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	if client.Type == "SPA" {
		//? set session to cookies
		dataSessions := map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		}

		err = utils.SetAndSaveSession(dataSessions, sess, ctx.Request(), ctx.Response())
		if err != nil {
			errorMessage := fmt.Sprintf("failed to set session with error: %v", err)
			utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

			utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

			return nil
		}
	}

	authToken := entity.Token{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
		ExpiresIn:    int(tokenExpTime),
		Scope:        auth.Scope.String,
		TokenType:    "Bearer",
		RedirectUri:  params.RedirectUri,
	}

	//? update auth use date so it didn't used twice
	err = tokenService.UpdateAuth(context, params.Code)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to update auth code with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? insert refresh token to db
	err = tokenService.CreateAuth(context, *refreshToken, auth.Scope.String, auth.UserID.String, 2)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to create refresh token with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	return ctx.JSON(http.StatusOK, utils.GenerateResponseJson(nil, true, authToken))
}
