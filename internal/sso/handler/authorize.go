package handler

import (
	"app/api/sso/operation"
	"app/internal/sso/entity"
	"app/internal/sso/logic"
	"app/internal/sso/repository"
	"app/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (r *RestService) Authorize(ctx echo.Context, params *operation.AuthorizeRequest) error {
	context := context.Background()
	serviceName := "GET Authorize"

	repo := repository.InitRepo(r.dbr, r.dbw)
	lc := logic.InitLogic()
	authorizeService := repository.AuthorizeRepository(repo)
	authorizeLogic := logic.AuthorizeLogic(lc)

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

		//TODO: ganti ke redirect auth server login page
		return ctx.NoContent(http.StatusNotImplemented)
	}

	//? if guid found, then current request can skip login and continue process auth code
	guid := sess.Values["guid"]

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

		//TODO: redirect ke auth server login page karena sesi tidak ketemu
		return ctx.NoContent(http.StatusNotImplemented)
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

	redirectUrls, err := authorizeService.FetchClientRedirects(context, params.ClientId)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to fetch redirect urls with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	redirectValid := authorizeLogic.ValidateRedirectUrls(redirectUrls, params.RedirectUrl)
	if !redirectValid {
		errorMessage := "redirect url invalid"
		utils.WarningLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusForbidden, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//TODO: coba cek masih perlu validasi lain ga, klo ga ada langsung simpen auth code nya ke db dan lanjut ke /token dengan flow jika auth code ditemukan langsung hapus

	authorizeCode, err := utils.GenerateRandomString(64)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to generate auth code with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	res := entity.AuthorizeResponse{
		AuthorizeCode: *authorizeCode,
	}

	return ctx.JSON(http.StatusOK, utils.GenerateResponseJson(nil, true, res))
}
