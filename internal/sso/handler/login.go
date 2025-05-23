package handler

import (
	"app/api/sso/operation"
	"app/internal/sso/entity"
	"app/internal/sso/repository"
	"app/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (r *RestService) Login(ctx echo.Context, params *operation.LoginRequest) error {
	context := context.Background()

	repo := repository.InitRepo(r.dbr, r.dbw)
	loginService := repository.LoginRepository(repo)

	user, err := loginService.GetUser(context, params.Username)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to get user with error: %v", err)
		logrus.Warn(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}
	if user == nil {
		errorMessage := "username or password is wrong, please try again."
		logrus.Info(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusForbidden, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? hash password inputed
	inputPassHash, err := utils.HashBcrypt(params.Password)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to hash password with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? compare hash password
	if utils.ValidateHash(user.Password, *inputPassHash) {
		errorMessage := "username or password is wrong, please try again."
		logrus.Info(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusForbidden, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	sessionId, err := ctx.Cookie("Session-Id")
	if err != nil {
		errorMessage := "session id not found, please login again later"
		logrus.Warn(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusForbidden, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	session, err := loginService.GetSession(context, sessionId.Value)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to get session from db with error: %v", err)
		logrus.Warn(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}
	if session == nil {
		errorMessage := "session not found in db, please try again later"
		logrus.Info(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusForbidden, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	err = loginService.UpdateUserIdSession(context, params.Username, sessionId.Value)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to update session in database with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	authorizationCode, err := utils.GenerateRandomString(64)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to generate auth code with error: %v", err)
		logrus.Warn(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	err = loginService.CreateAuthToken(context, *authorizationCode, sessionId.Value)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to insert auth token with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	message := entity.LoginEncrypt{
		CodeChallenge:       session.CodeChallenge,
		CodeChallengeMethod: session.CodeChallengeMethod,
		AuthorizationCode:   *authorizationCode,
	}

	messageString, _ := json.Marshal(message)

	ciphertext, err := utils.EncryptJwe(string(messageString), session.ClientID)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to encrypt message with error: %v", err)
		logrus.Warn(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	callbackUrl := viper.GetString(fmt.Sprintf("secret.%v.callback_url", session.ClientID))
	redParams := url.Values{}
	redParams.Add("code", *ciphertext)

	return ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("%v?%v", callbackUrl, redParams.Encode()))
}
