package main

import (
	"fmt"
	"log"
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
	defer cfg.wg.Done()

	// normalize currentURL
	currentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("error after normalization, %v\n", err)
		return
	}

	// check is current URL contains Base URL, does not add to apge if false
	// if !strings.Contains(currentURL, cfg.baseURL.Host) {
	// 	return
	// }

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
		normUrl, err := normalizeURL(URL)
		if err != nil {
			log.Fatal(err)
		}
		// check is URL contains Base URL, does not add to pages map if false
		if !strings.Contains(normUrl, cfg.baseURL.Host) {
			continue
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
	cfg.pages[normalizedURL] += 1
	if cfg.pages[normalizedURL] == 1 {
		cfg.mu.Unlock()
		return true
	}
	cfg.mu.Unlock()
	return false
}
