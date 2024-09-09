package main

import (
	"fmt"
	"os"
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

		pages := make(map[string]int)
		baseURL, err := normalizeURL(args[0])
		if err != nil {
			fmt.Println(err)
		}
		pages[baseURL] = 0

		crawlPage(args[0], args[0], pages)
		for k, v := range pages {
			fmt.Printf("Key: %s Val: %v\n", k, v)
		}
		fmt.Println("Ending Crawl")
	}
}
