package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	args := os.Args[1:]

	switch {
	case len(args) < 1:
		fmt.Println("no website provided")
		os.Exit(1)
	case len(args) > 3:
		fmt.Println("too many arguments provided")
		os.Exit(1)
	default:
		fmt.Printf("starting crawl of %s\n", args[0])

		var cfg config
		cfg.pages = make(map[string]int)

		baseURL, err := url.Parse(args[0])
		if err != nil {
			fmt.Println(err)
		}
		cfg.baseURL = baseURL

		if len(args) == 2 {
			buffer, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cfg.concurrencyControl = make(chan struct{}, buffer)
		} else {
			cfg.concurrencyControl = make(chan struct{}, 5)
		}

		if len(args) == 3 {
			maxPages, err := strconv.Atoi(args[2])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cfg.maxPages = maxPages
		} else {
			cfg.maxPages = 100
		}
		cfg.mu = &sync.Mutex{}
		cfg.wg = &sync.WaitGroup{}

		cfg.wg.Add(1)
		go cfg.crawlPage(args[0])
		cfg.wg.Wait()
		printReport(cfg.pages, args[0])
	}
}
