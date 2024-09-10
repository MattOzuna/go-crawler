package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {

	body, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		fmt.Println("parse fail")
		return []string{}, err
	}

	parsedURLs := []string{}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					parsedURLs = append(parsedURLs, a.Val)
					break
				}
			}
		}
		for node := n.FirstChild; node != nil; node = node.NextSibling {
			f(node)
		}
	}
	f(body)

	for i, url := range parsedURLs {
		if !strings.HasPrefix(url, "/") {
			url = "/" + url
		}
		if !strings.Contains(url, "http") {
			parsedURLs[i] = rawBaseURL + url
		}
	}

	return parsedURLs, nil
}
