package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Problem struct {
	Type     string      `json:"type"`
	Title    string      `json:"title"`
	Status   int         `json:"status"`
	Detail   interface{} `json:"detail"`
	Instance string      `json:"instance"`
	Guid     string      `json:"guid"`
}

func generateProblemJson(statusCode int, message, instance, guid string) Problem {
	return Problem{
		Title:    http.StatusText(statusCode),
		Status:   statusCode,
		Detail:   message,
		Instance: instance,
		Guid:     guid,
	}
}

func SendProblemDetailJson(ctx echo.Context, statusCode int, message, instance, guid string) error {
	problem := generateProblemJson(statusCode, message, instance, guid)

	ctx.Response().Header().Set("Content-Type", "application/problem+json")
	return ctx.JSON(statusCode, problem)
}
