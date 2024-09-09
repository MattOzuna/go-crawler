package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	// normalize currentURL
	currentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// normalize baseURL
	baseURL, err := normalizeURL(rawBaseURL)
	if err != nil {
		fmt.Println(err)
	}
	// check is current URL contains Base URL, and exits early if it does not
	if !strings.Contains(currentURL, baseURL) {
		return
	}

	// base case
	if pages[currentURL] == 1 {
		return
	}

	// get the HTML from currentURL and grab the URLs in <a> tags
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	URLs, err := getURLsFromHTML(html, rawBaseURL)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// mark that we have visited the current page in pages map
	// fmt.Printf("Visited %s\n", rawCurrentURL)
	pages[currentURL] = 1

	// add new URLs into pages map and crawl through them
	for _, url := range URLs {
		normUrl, err := normalizeURL(url)
		if err != nil {
			log.Fatal(err)
		}

		_, ok := pages[normUrl]
		if !ok {
			pages[normUrl] = 0
			crawlPage(rawBaseURL, url, pages)
		}
	}
}
