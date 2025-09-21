//go:generate stringer -type=ErrorCode
package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ErrorCode int

const (
	ERR_UNAUTHENTICATED ErrorCode = 401 //ERR_UNAUTHENTICATED
	ERR_INTERNAL_ERROR  ErrorCode = 500 //ERR_INTERNAL_ERROR
	ERR_FORBIDDEN       ErrorCode = 403 //ERR_FORBIDDEN
	ERR_BAD_REQUEST     ErrorCode = 400 //ERR_BAD_REQUEST
	ERR_CONFLICT        ErrorCode = 409 //ERR_CONFLICT
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

func generateProblemJson(statusCode ErrorCode, message, instance, guid string) Problem {
	return Problem{
		Type:     statusCode.String(),
		Title:    http.StatusText(int(statusCode)),
		Status:   int(statusCode),
		Message:  message,
		Instance: instance,
		Guid:     guid,
	}
}

func SendProblemDetailJson(ctx echo.Context, statusCode ErrorCode, message, instance, guid string) error {
	problem := generateProblemJson(statusCode, message, instance, guid)

	ctx.Response().Header().Set("Content-Type", "application/problem+json")
	return ctx.JSON(int(statusCode), problem)
}

func SendProblemDetailJsonValidate(ctx echo.Context, statusCode ErrorCode, message, instance, guid string, errors validator.ValidationErrors) error {
	errorKv := map[string]string{}

	for _, v := range errors {
		ns := v.Namespace()
		keys := strings.Split(ns, ".")

		errorKv[keys[1]] = normalizeError(v)
	}

	problem := generateProblemJson(statusCode, message, instance, guid)
	problem.Errors = errorKv

	ctx.Response().Header().Set("Content-Type", "application/problem+json")
	return ctx.JSON(int(statusCode), problem)
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

func ErrorCodeString(code string) string {
	switch code {
	case "ERR_UNAUTHENTICATED":
		return "current request can't be process because request doesn't meet required information to be proceed, please try to login again and try again later, if error still persist please contact our admin."
	case "ERR_INTERNAL_ERROR":
		return "current request can't be process because there is some error found when process, please contact our admin."
	case "ERR_FORBIDDEN":
		return "current request can't be process because one or more information to process current request is missing or not valid, please try to login again and try again later, if error still persist please contact our admin."
	case "ERR_BAD_REQUEST":
		return "current request can't be process because one or more of the requested properties is not valid, please check the request and try again."
	case "ERR_CONFLICT":
		return "current request can't be process because the request is conflicting with current process, please clear your cache and try again later, if error still persist please contact our admin."
	default:
		return fmt.Sprintf("error code of %v is currently not defined, please contact our admin for more information.", code)
	}
}
