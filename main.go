package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"
)

func main() {
	args := os.Args[1:]

	switch {
	case len(args) < 1:
		fmt.Println("no website provided")
		os.Exit(1)
	case len(args) > 1:
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
		cfg.pages[baseURL.Host] = 0
		cfg.baseURL = baseURL
		cfg.concurrencyControl = make(chan struct{}, 5)
		cfg.mu = &sync.Mutex{}
		cfg.wg = &sync.WaitGroup{}

		cfg.wg.Add(1)
		cfg.concurrencyControl <- struct{}{}
		go cfg.crawlPage(args[0])
		<-cfg.concurrencyControl
		cfg.wg.Wait()

		// for k, v := range pages {
		// 	fmt.Printf("Key: %s Val: %v\n", k, v)
		// }
		fmt.Println("Ending Crawl")
	}
}
