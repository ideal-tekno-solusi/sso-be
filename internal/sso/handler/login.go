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
	result := entity.Token{}

	user, err := loginService.GetUser(context, params.Username)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to get user with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}
	if user == nil {
		errorMessage := "username or password is wrong, please try again."
		logrus.Warn(errorMessage)

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
		logrus.Warn(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusForbidden, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	err = loginService.CreateSession(context, params.State, params.Username, params.ClientId, params.CodeChallenge, params.CodeChallengeMethod, params.Scopes, params.RedirectUrl)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to create session with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? req GET to authorize
	authorizeDomain := viper.GetString("config.url.internal.domain")
	authorizePath := viper.GetString("config.url.internal.path.authorize")
	codeVerifier, err := ctx.Cookie("verifier")
	if err != nil {
		errorMessage := "failed to get code verifier, please try login again"
		logrus.Warn(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusUnauthorized, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	verifierAge := viper.GetInt("config.verifier.age")
	verifierDomain := viper.GetString("config.verifier.domain")
	verifierPath := viper.GetString("config.verifier.path")
	verifierSecure := viper.GetBool("config.verifier.secure")
	verifierHttponly := viper.GetBool("config.verifier.httponly")

	cookies := []*http.Cookie{
		{Name: codeVerifier.Name, Value: codeVerifier.Value, Path: verifierPath, Domain: verifierDomain, MaxAge: verifierAge, Secure: verifierSecure, HttpOnly: verifierHttponly},
	}

	query := url.Values{}
	query.Add("response_type", params.ResponseType)
	query.Add("client_id", params.ClientId)
	query.Add("redirect_url", params.RedirectUrl)
	query.Add("scopes", params.Scopes)
	query.Add("state", params.State)
	query.Add("code_challenge", params.CodeChallenge)
	query.Add("code_challenge_method", params.CodeChallengeMethod)

	status, res, err := utils.SendHttpGetRequest(fmt.Sprintf("%v%v", authorizeDomain, authorizePath), &query, cookies)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to request authorize with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}
	if status != http.StatusOK {
		errorMessage := fmt.Sprintf("response from server is not ok, status %v", status)
		logrus.Warn(errorMessage)

		utils.SendProblemDetailJson(ctx, status, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	err = json.Unmarshal(res, &result)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to unmarshal message with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	return ctx.JSON(http.StatusOK, result)
}
