package handler

import (
	"app/api/sso/operation"
	"app/internal/sso/repository"
	"app/utils"
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
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

		//TODO: should return message too
		return ctx.NoContent(http.StatusForbidden)
	}

	//TODO: hash plain password from params and compare with user data from db
	if params.Password != user.Password {
		errorMessage := "username or password is wrong, please try again."
		logrus.Info(errorMessage)

	}

	return ctx.NoContent(http.StatusOK)
}
