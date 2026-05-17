package linescaler_test

import (
	"fmt"

	"github.com/yourorg/logslice/internal/linescaler"
)

// ExampleScaler_Scale demonstrates converting a millisecond latency field
// (column 2, space-delimited) to seconds by scaling by 0.001.
func ExampleScaler_Scale() {
	s := linescaler.New(0.001, 2, " ")
	out, err := s.Scale("2024-01-15 INFO 1500 request completed")
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
	// Output: 2024-01-15 INFO 1.5 request completed
}

// ExampleScaler_Scale_disabled shows that a factor of 1 is a no-op.
func ExampleScaler_Scale_disabled() {
	s := linescaler.New(1, 0, " ")
	line := "42 unchanged line"
	out, _ := s.Scale(line)
	fmt.Println(out)
	// Output: 42 unchanged line
}
