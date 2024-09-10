package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	defer cfg.wg.Done()

	// normalize currentURL
	currentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// check is current URL contains Base URL, and return early if it does not
	if !strings.Contains(currentURL, cfg.baseURL.Host) {
		return
	}

	// get the HTML from currentURL and grab the URLs in <a> tags
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	URLs, err := getURLsFromHTML(html, cfg.baseURL.Host)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// base case
	if !cfg.addPageVisit(currentURL) {
		return
	}

	// add new URLs into pages map and crawl through them
	for _, URL := range URLs {
		normUrl, err := normalizeURL(URL)
		if err != nil {
			log.Fatal(err)
		}

		cfg.mu.Lock()
		_, ok := cfg.pages[normUrl]
		cfg.mu.Unlock()

		if !ok {
			cfg.mu.Lock()
			cfg.pages[normUrl] = 0
			cfg.mu.Unlock()

			cfg.wg.Add(1)
			cfg.concurrencyControl <- struct{}{}
			go cfg.crawlPage("https://" + URL)
			<-cfg.concurrencyControl

		}
	}
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	if cfg.pages[normalizedURL] == 1 {
		cfg.mu.Unlock()
		return false
	}

	// mark that we have visited the current page in pages map
	fmt.Printf("Visited %s\n", normalizedURL)
	cfg.pages[normalizedURL] = 1
	cfg.mu.Unlock()
	return true
}
