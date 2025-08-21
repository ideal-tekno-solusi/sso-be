package handler

import (
	"app/api/sso/operation"
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
	//TODO: lanjutin token untuk yg tuker auth code dan refresh token
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

	codeChallengeSource := sess.Values["code_Challenge"]
	codeChallengeMethodSource := sess.Values["code_challenge_method"]
	guid := sess.Values["guid"]
	if codeChallengeSource == nil || codeChallengeMethodSource == nil || guid == nil {
		errorMessage := "session is empty, please cleare cache and try to login again"
		utils.WarningLog(errorMessage, ctx.Path(), serviceName)

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
	if client.Type.String == "SS" {
		if client.Secret.String != params.ClientSecret {
			errorMessage := "client id or secret not valid"
			utils.WarningLog(errorMessage, ctx.Path(), serviceName)

			utils.SendProblemDetailJson(ctx, http.StatusForbidden, errorMessage, ctx.Path(), uuid.NewString())

			return nil
		}
	}

	// //? generate code challenge from code verifier req
	var codeChallenge string = codeChallengeSource.(string)

	if codeChallengeMethodSource == "SHA256" {
		hash := sha256.New()
		hash.Write([]byte(params.CodeVerifier))
		codeChallenge = base64.StdEncoding.EncodeToString(hash.Sum(nil))
	}

	//? check if code verifier equal to code challenge
	if params.CodeVerifier != codeChallenge {
		errorMessage := "code verifier not match"
		utils.WarningLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusForbidden, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//TODO: lanjut ambil auth code nya di table auth untuk cek apakah auth code sudah digunakan atau belum, lalu update kolom use_date dan generate refresh token lalu masukkan ke table auth sebagai code baru

	// //? get token by code challenge and check validity of code verifier at once
	// token, err := tokenService.GetToken(context, codeChallenge)
	// if err != nil {
	// 	errorMessage := fmt.Sprintf("failed to get token with error: %v", err)
	// 	logrus.Error(errorMessage)

	// 	utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

	// 	return nil
	// }

	// //? check if auth code is match
	// if params.Code != token.ID {
	// 	errorMessage := "token not match, please login again later"
	// 	logrus.Warn(errorMessage)

	// 	utils.SendProblemDetailJson(ctx, http.StatusUnauthorized, errorMessage, ctx.Path(), uuid.NewString())

	// 	return nil
	// }

	// //? generate jwt token
	// jwtBody := map[string]string{
	// 	"username":    token.Username,
	// 	"name":        token.Name,
	// 	"redirectUrl": token.RedirectUrl.String,
	// }

	// tokenExpTime := viper.GetInt("secret.expToken")
	// refreshExpTime := viper.GetInt("secret.refreshToken")

	// accessToken, err := utils.GenerateAuthToken(jwtBody, tokenExpTime)
	// if err != nil {
	// 	errorMessage := fmt.Sprintf("failed to generate access token with error: %v", err)
	// 	logrus.Error(errorMessage)

	// 	utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

	// 	return nil
	// }

	// refreshToken, err := utils.GenerateAuthToken(jwtBody, refreshExpTime)
	// if err != nil {
	// 	errorMessage := fmt.Sprintf("failed to generate access token with error: %v", err)
	// 	logrus.Error(errorMessage)

	// 	utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

	// 	return nil
	// }

	// err = tokenService.CreateRefreshToken(context, *refreshToken, token.Username)
	// if err != nil {
	// 	errorMessage := fmt.Sprintf("failed to insert refresh token with error: %v", err)
	// 	logrus.Error(errorMessage)

	// 	utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

	// 	return nil
	// }

	// authToken := entity.Token{
	// 	AccessToken:  *accessToken,
	// 	RefreshToken: *refreshToken,
	// 	ExpiresIn:    tokenExpTime,
	// 	Scope:        token.Scopes.String,
	// 	TokenType:    "Bearer",
	// }

	// //? delete auth token and session and remove session id
	// err = tokenService.DeleteAuthToken(context, token.SessionID.String)
	// if err != nil {
	// 	errorMessage := fmt.Sprintf("failed to delete auth token with error: %v", err)
	// 	logrus.Error(errorMessage)

	// 	utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

	// 	return nil
	// }

	// err = tokenService.DeleteSession(context, token.SessionID.String)
	// if err != nil {
	// 	errorMessage := fmt.Sprintf("failed to delete session with error: %v", err)
	// 	logrus.Error(errorMessage)

	// 	utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

	// 	return nil
	// }

	// return ctx.JSON(http.StatusOK, utils.GenerateResponseJson(nil, true, authToken))
	return ctx.NoContent(http.StatusNotImplemented)
}
