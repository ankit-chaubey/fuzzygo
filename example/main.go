package main

import (
	"fmt"
	"time"

	"github.com/ankit-chaubey/fuzzygo/fuzzy"
)

func main() {
	methods, err := fuzzy.LoadTLMethods("output.json")
	if err != nil {
		panic(err)
	}

	start := time.Now()
	results, total := fuzzy.Rank("mesages.getMesage", methods, 10)
	elapsed := time.Since(start)

	fmt.Printf("Found %d matches in %.3f ms\n",
		total,
		float64(elapsed.Microseconds())/1000,
	)

	for i, r := range results {
		fmt.Printf("%d. %s (score=%d)\n", i+1, r.Item, r.Score)
	}
}
