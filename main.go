package main

import (
	"fmt"
	"log"
	"mapreduce/worker"
	"os"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	start := time.Now()

	files, err := os.ReadDir("input")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		wg.Add(1)
		k := fmt.Sprintf("input/%s", file.Name())
		go worker.Map(k, &wg)
	}
	wg.Wait()

	worker.Reduce()

	elapsed := time.Since(start)
	fmt.Printf("%s\n", elapsed)
}
