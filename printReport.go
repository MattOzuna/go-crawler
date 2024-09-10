package main

import (
	"fmt"
	"slices"
)

func printReport(pages map[string]int, baseURL string) {
	fmt.Println("========================================")
	fmt.Printf("REPORT for %s\n", baseURL)
	fmt.Println("========================================")

	keys := make([]string, len(pages))
	values := make([]int, len(pages))
	i := 0
	for k, v := range pages {
		keys[i] = k
		values[i] = v
		i++
	}

	slices.Sort(keys)
	slices.Sort(values)
	slices.Reverse(values)

	for i = 0; i < len(pages); i++ {
		fmt.Printf("Found %v internal links to %s\n", values[i], keys[i])
	}

}
