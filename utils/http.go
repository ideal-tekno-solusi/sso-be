package utils

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

func SendHttpPostRequest(url string, body []byte, cookies []*http.Cookie) (int, []byte, error) {
	method := "POST"

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	//TODO: keknya bakal perlu penjagaan lebih deh ini
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return res.StatusCode, resBody, nil
}

func SendHttpGetRequest(url string, body *url.Values, cookies []*http.Cookie) (int, []byte, error) {
	method := "GET"

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if body != nil {
		req.URL.RawQuery = body.Encode()
	}

	req.Header.Add("Content-Type", "application/json")

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	//TODO: keknya bakal perlu penjagaan lebih deh ini
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return res.StatusCode, resBody, nil
}
