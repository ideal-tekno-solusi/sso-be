package handler

import (
	"app/api/sso/operation"
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

func (r *RestService) Login(ctx echo.Context, params *operation.LoginRequest) error {
	context := context.Background()
	serviceName := "POST Login"

	repo := repository.InitRepo(r.dbr, r.dbw)
	loginService := repository.LoginRepository(repo)
	guid := uuid.NewString()
	req := operation.AuthorizeRequest{}

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

	//? send error if already login
	alreadyLogin := sess.Values["guid"]
	if alreadyLogin != nil {
		errorMessage := "failed to process login because user already login"
		utils.WarningLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusBadRequest, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	user, err := loginService.GetUser(context, params.Username)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to get user with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}
	if user == nil {
		errorMessage := "username or password is wrong, please try again."
		utils.WarningLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusForbidden, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? hash password inputed
	inputPassHash, err := utils.HashBcrypt(params.Password)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to hash password with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? compare hash password
	if utils.ValidateHash(user.Password, *inputPassHash) {
		errorMessage := "username or password is wrong, please try again."
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusForbidden, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? get request query from session
	reqByte := sess.Values["authorization"]
	if reqByte == nil {
		errorMessage := "authorization property not found in current session, please try to request authorization again later"
		utils.WarningLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusUnauthorized, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	err = json.Unmarshal(reqByte.([]byte), &req)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to unmarshal request with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? save guid that be used for session to db
	err = loginService.CreateSession(context, guid, user.ID)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to add session to db with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? set session to cookies
	dataSessions := map[string]interface{}{
		"guid":                  guid,
		"code_Challenge":        req.CodeChallenge,
		"code_challenge_method": req.CodeChallengeMethod,
	}

	err = utils.SetAndSaveSession(dataSessions, sess, ctx.Request(), ctx.Response())
	if err != nil {
		errorMessage := fmt.Sprintf("failed to set session with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	//? set up redirect to /authorize
	authDomain := viper.GetString("config.url.internal.domain")
	authPath := viper.GetString("config.url.internal.path.authorize")

	query := url.Values{}
	query.Add("response_type", req.ResponseType)
	query.Add("client_id", req.ClientId)
	query.Add("redirect_uri", req.RedirectUri)
	query.Add("scope", req.Scope)
	query.Add("state", req.State)
	query.Add("code_challenge", req.CodeChallenge)
	query.Add("code_challenge_method", req.CodeChallengeMethod)

	u, err := url.Parse(fmt.Sprintf("%v%v", authDomain, authPath))
	if err != nil {
		errorMessage := fmt.Sprintf("failed to parse url with error: %v", err)
		utils.ErrorLog(errorMessage, ctx.Path(), serviceName)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	u.RawQuery = query.Encode()

	//? delete authorization from session
	deleteSession := []string{
		"authorization",
	}

	utils.DeleteAndSaveSession(deleteSession, sess, ctx.Request(), ctx.Response())

	return ctx.Redirect(http.StatusFound, u.String())
}
