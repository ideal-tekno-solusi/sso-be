package handler

import (
	"app/api/sso/operation"
	"app/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (r *RestService) User(ctx echo.Context, params *operation.UserRequest) error {
	token := strings.Split(params.Token, " ")

	plainText, err := utils.DecryptJwt(token[1])
	if err != nil {
		errorMessage := fmt.Sprintf("failed to decrypt jwt with error: %v", err)
		logrus.Error(errorMessage)

		utils.SendProblemDetailJson(ctx, http.StatusInternalServerError, errorMessage, ctx.Path(), uuid.NewString())

		return nil
	}

	return ctx.JSON(http.StatusOK, utils.GenerateResponseJson(true, plainText))
}
