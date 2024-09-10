package main

import (
	"net/url"
	"strings"
)

func normalizeURL(URL string) (string, error) {
	URL = strings.TrimRight(URL, "/")
	URL = strings.ToLower(URL)
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return "", err
	}

	return parsedURL.Host + parsedURL.Path, nil
}
