package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	res, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode > 399 {
		return "", fmt.Errorf("unsuccessful request status code: %v", res.StatusCode)
	}

	contentHeader := res.Header.Get("Content-type")
	if strings.Contains(contentHeader, "type/html") {
		return "", errors.New("content type not text/html")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), err
}
