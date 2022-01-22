package utils

import (
	"net/http"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   error
}

type ErrorResp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func HttpGet(url string) (*http.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return res, nil
}
