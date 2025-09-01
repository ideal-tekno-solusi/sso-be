package utils

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func SetAndSaveSession(data map[string]interface{}, sess *sessions.Session, req *http.Request, res *echo.Response) error {
	for k, v := range data {
		sess.Values[k] = v
	}

	err := sess.Save(req, res)
	if err != nil {
		return err
	}

	return nil
}

func DeleteAndSaveSession(data []string, sess *sessions.Session, req *http.Request, res *echo.Response) error {
	for _, v := range data {
		delete(sess.Values, v)
	}

	err := sess.Save(req, res)
	if err != nil {
		return err
	}

	return nil
}

func DeleteSession(sess *sessions.Session, req *http.Request, res *echo.Response) error {
	sess.Options.MaxAge = -1

	err := sess.Save(req, res)
	if err != nil {
		return err
	}

	return nil
}
