package handler

import (
	"app/api/sso/operation"
	"app/internal/sso/logic"
	"app/internal/sso/repository"
	"app/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func (r *RestService) Authorize(ctx echo.Context, params *operation.AuthorizeRequest) error {
	context := context.Background()
	serviceName := "GET Authorize"

	repo := repository.InitRepo(r.dbr, r.dbw)
	lc := logic.InitLogic()
	authorizeService := repository.AuthorizeRepository(repo)
	authorizeLogic := logic.AuthorizeLogic(lc)
	redirectFe := viper.GetString("config.url.redirect_fe.login")

	//? get session, will create new if not found
	sess, err := session.Get("session", ctx)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to get session with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? save requested param as string and redirect to auth login page
	if sess.IsNew {
		paramString, err := json.Marshal(params)
		if err != nil {
			errorMessage := fmt.Sprintf("failed to marshal params with error: %v", err)
			utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

			utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

			return nil
		}

		//? set session to cookies
		dataSessions := map[string]interface{}{
			"authorization": paramString,
		}

		err = utils.SetAndSaveSession(dataSessions, sess, ctx.Request(), ctx.Response())
		if err != nil {
			errorMessage := fmt.Sprintf("failed to set session with error: %v", err)
			utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

			utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

			return nil
		}

		return ctx.Redirect(http.StatusFound, redirectFe)
	}

	//? if guid found, then current request can skip login and continue process auth code
	guid := sess.Values["guid"]
	if guid == nil {
		errorMessage := "session not found, please clear cache and try to login again"
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	existSession, err := authorizeService.GetSession(context, guid.(string))
	if err != nil {
		errorMessage := fmt.Sprintf("failed to get existing session with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}
	if existSession == nil {
		errorMessage := "session not exist, please login again later"
		utils.WarningLog(errorMessage, ctx.Path(), serviceName)

		return ctx.Redirect(http.StatusFound, redirectFe)
	}

	client, err := authorizeService.GetClient(context, params.ClientId)
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

	redirectUris, err := authorizeService.FetchClientRedirects(context, params.ClientId)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to fetch redirect urls with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	redirectValid := authorizeLogic.ValidateRedirectUris(redirectUris, params.RedirectUri)
	if !redirectValid {
		errorMessage := "redirect url invalid"
		utils.WarningLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusForbidden, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	authorizeCode, err := utils.GenerateRandomString(64)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to generate auth code with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? save auth code and user id to db
	err = authorizeService.CreateAuth(context, *authorizeCode, params.Scope, existSession.UserID.String, 1)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to create auth with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? set up redirect to /callback fe
	query := url.Values{}
	query.Add("code", *authorizeCode)

	u, err := url.Parse(params.RedirectUri)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to parse url with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	u.RawQuery = query.Encode()

	return ctx.Redirect(http.StatusFound, u.String())
}
