package main

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	// normalize currentURL
	currentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("error after normalization, %v\n", err)
		return
	}

	// check is current URL contains Base URL, does not add to apge if false
	if !strings.Contains(currentURL, cfg.baseURL.Host) {
		return
	}

	// base case: return if page has been visited
	isFirst := cfg.addPageVisit(currentURL)
	if !isFirst {
		return
	}

	// return if max pages visited is reached
	cfg.mu.Lock()
	currPageLen := len(cfg.pages)
	cfg.mu.Unlock()
	if currPageLen > cfg.maxPages {
		return
	}

	// get the HTML from currentURL and grab the URLs in <a> tags
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
	}
	URLs, err := getURLsFromHTML(html, cfg.baseURL.Host)
	if err != nil {
		fmt.Println(err)
	}

	// add new URLs into pages map and crawl through them
	for _, URL := range URLs {

		cfg.wg.Add(1)
		go cfg.crawlPage("https://" + URL)

	}
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, visited := cfg.pages[normalizedURL]; visited {
		cfg.pages[normalizedURL]++
		return false
	}

	cfg.pages[normalizedURL] = 1
	return true
}
