package linethrottler_test

import (
	"fmt"
	"time"

	"github.com/yourorg/logslice/internal/linethrottler"
)

func ExampleThrottler_Apply() {
	// Disabled throttler — no delay, lines pass through immediately.
	th := linethrottler.New(0)
	lines := []string{"2024-01-01 info start", "2024-01-01 info stop"}
	for _, l := range lines {
		fmt.Println(th.Apply(l))
	}
	// Output:
	// 2024-01-01 info start
	// 2024-01-01 info stop
}

func ExampleThrottler_Apply_enabled() {
	// Enabled throttler — sleeps between lines.
	// In production this would use a real delay; here we show the API.
	th := linethrottler.New(10 * time.Millisecond)
	line := th.Apply("2024-01-01 warn something happened")
	fmt.Println(line)
	// Output:
	// 2024-01-01 warn something happened
}
