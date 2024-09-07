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
		fmt.Printf("starting crawl of %s", args[0])
		htmlRes, err := getHTML(args[0])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print(htmlRes)
	}
}
