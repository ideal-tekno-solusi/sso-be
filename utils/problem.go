package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Problem struct {
	Type     string      `json:"type"`
	Title    string      `json:"title"`
	Status   int         `json:"status"`
	Detail   interface{} `json:"detail,omitempty"`
	Message  string      `json:"message"`
	Errors   interface{} `json:"errors,omitempty"`
	Instance string      `json:"instance"`
	Guid     string      `json:"guid"`
}

func generateProblemJson(statusCode int, message, instance, guid string) Problem {
	return Problem{
		Type:     "about:blank",
		Title:    http.StatusText(statusCode),
		Status:   statusCode,
		Message:  message,
		Instance: instance,
		Guid:     guid,
	}
}

func SendProblemDetailJson(ctx echo.Context, statusCode int, message, instance, guid string) error {
	problem := generateProblemJson(statusCode, message, instance, guid)

	ctx.Response().Header().Set("Content-Type", "application/problem+json")
	return ctx.JSON(statusCode, problem)
}

func SendProblemDetailJsonValidate(ctx echo.Context, statusCode int, message, instance, guid string, errors validator.ValidationErrors) error {
	errorKv := map[string]string{}

	for _, v := range errors {
		ns := v.Namespace()
		keys := strings.Split(ns, ".")

		errorKv[keys[1]] = normalizeError(v)
	}

	problem := generateProblemJson(statusCode, message, instance, guid)
	problem.Errors = errorKv

	ctx.Response().Header().Set("Content-Type", "application/problem+json")
	return ctx.JSON(statusCode, problem)
}

func normalizeError(err validator.FieldError) string {
	//? every time using new validator, need to define it's error message here
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%v is required", strings.Split(err.Namespace(), ".")[1])
	case "max":
		return fmt.Sprintf("%v max length is %v", strings.Split(err.Namespace(), ".")[1], err.Param())
	case "eq":
		return fmt.Sprintf("%v only accept %v", strings.Split(err.Namespace(), ".")[1], err.Param())
	case "oneofci":
		return fmt.Sprintf("%v contain unidentified string", strings.Split(err.Namespace(), ".")[1])
	default:
		return "undefined error"
	}
}
