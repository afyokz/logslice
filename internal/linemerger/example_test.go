package linemerger_test

import (
	"fmt"

	"github.com/yourorg/logslice/internal/linemerger"
)

func ExampleMerger_Feed() {
	lines := []string{
		"2024-01-15 12:00:00 user login uid=42",
		"2024-01-15 12:00:00 user login detail: success",
		"2024-01-15 12:00:01 user logout uid=42",
	}

	m := linemerger.New(19, " | ")
	for _, l := range lines {
		if out, ok := m.Feed(l); ok {
			fmt.Println(out)
		}
	}
	if out, ok := m.Flush(); ok {
		fmt.Println(out)
	}

	// Output:
	// 2024-01-15 12:00:00 user login uid=42 | 2024-01-15 12:00:00 user login detail: success
	// 2024-01-15 12:00:01 user logout uid=42
}

func ExampleMerger_Flush() {
	m := linemerger.New(10, "-")
	m.Feed("2024-01-01 first")
	m.Feed("2024-01-01 second")

	out, ok := m.Flush()
	fmt.Println(ok, out)

	// Output:
	// true 2024-01-01 first-2024-01-01 second
}
