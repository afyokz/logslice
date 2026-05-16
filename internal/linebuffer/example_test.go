package linebuffer_test

import (
	"fmt"
	"strings"

	"github.com/logslice/logslice/internal/linebuffer"
)

// ExampleBuffer_Lines demonstrates capturing trailing context around a match.
func ExampleBuffer_Lines() {
	buf := linebuffer.New(3)

	lines := []string{
		"2024-01-01 DEBUG init",
		"2024-01-01 INFO  starting",
		"2024-01-01 ERROR something failed",
		"2024-01-01 DEBUG cleanup",
	}

	for _, line := range lines {
		buf.Push(line)
		if strings.Contains(line, "ERROR") {
			for _, ctx := range buf.Lines() {
				fmt.Println(ctx)
			}
		}
	}
	// Output:
	// 2024-01-01 DEBUG init
	// 2024-01-01 INFO  starting
	// 2024-01-01 ERROR something failed
}

// ExampleBuffer_Reset demonstrates clearing the buffer between matches.
func ExampleBuffer_Reset() {
	buf := linebuffer.New(4)
	buf.Push("line1")
	buf.Push("line2")
	fmt.Println("before reset:", buf.Len())
	buf.Reset()
	fmt.Println("after reset:", buf.Len())
	// Output:
	// before reset: 2
	// after reset: 0
}
