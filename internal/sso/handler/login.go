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
	authorizeRes := entity.AuthorizeResponse{}
	response := entity.LoginResponse{}

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

	query := url.Values{}
	query.Add("responseType", params.ResponseType)
	query.Add("clientId", params.ClientId)
	query.Add("redirectUrl", params.RedirectUrl)
	query.Add("scopes", params.Scopes)
	query.Add("state", params.State)
	query.Add("codeChallenge", params.CodeChallenge)
	query.Add("codeChallengeMethod", params.CodeChallengeMethod)

	status, res, err := utils.SendHttpGetRequest(fmt.Sprintf("%v%v", authorizeDomain, authorizePath), &query, nil)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to request authorize with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}
	if status != http.StatusOK {
		errorMessage := fmt.Sprintf("response from server is not ok, response server: %v", string(res))
		logrus.Warn(errorMessage)

		utils.SendProblemDetailJson(ctx, status, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	reqDefaultRes, reqBodyRes, err := utils.BindResponse(res)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to bind response with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	err = json.Unmarshal(*reqBodyRes, &authorizeRes)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to unmarshal message with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	callbackUrl := viper.GetString(fmt.Sprintf("secret.%v.callback_url", params.ClientId))
	redParams := url.Values{}
	redParams.Add("code", authorizeRes.AuthorizeCode)
	redParams.Add("state", params.State)

	response.CallbackUrl = fmt.Sprintf("%v?%v", callbackUrl, redParams.Encode())

	return ctx.JSON(http.StatusOK, utils.GenerateResponseJson(reqDefaultRes, true, response))
}
