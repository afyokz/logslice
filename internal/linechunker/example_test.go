package linechunker_test

import (
	"fmt"

	"github.com/yourorg/logslice/internal/linechunker"
)

// ExampleChunker_Feed demonstrates batching lines into chunks of two.
func ExampleChunker_Feed() {
	c := linechunker.New(2)

	lines := []string{"alpha", "beta", "gamma", "delta", "epsilon"}
	var batches [][]string

	for _, l := range lines {
		if batch := c.Feed(l); batch != nil {
			batches = append(batches, batch)
		}
	}
	if rem := c.Flush(); rem != nil {
		batches = append(batches, rem)
	}

	for i, b := range batches {
		fmt.Printf("batch %d: %v\n", i+1, b)
	}
	// Output:
	// batch 1: [alpha beta]
	// batch 2: [gamma delta]
	// batch 3: [epsilon]
}

// ExampleChunker_Feed_disabled shows that a zero chunk size passes every
// line through immediately as a single-element slice.
func ExampleChunker_Feed_disabled() {
	c := linechunker.New(0)
	got := c.Feed("hello")
	fmt.Println(got)
	// Output:
	// [hello]
}
