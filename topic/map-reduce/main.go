package main

import (
	"fmt"
	"strings"
	"sync"
)

func mapF(docID string, contents string) []KeyValue {
	// Split the contents into individual words
	words := strings.Fields(contents)

	// Create a slice to hold the KeyValue pairs
	var kvs []KeyValue

	// Iterate over each word and create a KeyValue pair
	for _, w := range words {
		// Convert the word to lowercase and set the value to "1"
		kvs = append(kvs, KeyValue{
			Key:   strings.ToLower(w),
			Value: "1",
		})
	}

	// Return the slice of KeyValue pairs
	return kvs
}

func reduceF(key string, values []string) string {
	return fmt.Sprintf("%d", len(values))
}

func main() {
	inputs := []string{
		"foo bar baz foo",
		"bar baz qux",
		"foo qux qux qux",
		"lorem ipsum dolor sit amet",
	}

	master := NewMaster(inputs, 3)

	numWorkers := 4
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		w := NewWorker(i, master, mapF, reduceF)
		go w.Run(&wg)
	}

	wg.Wait()

	for i := 0; i < master.NumReduce; i++ {
		fmt.Printf("\n--- reducer %d output ---\n%s\n", i, master.Outputs[i])
	}
}
