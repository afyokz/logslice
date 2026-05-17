// Package linethrottler provides a Throttler that introduces a
// configurable inter-line delay during log output.
//
// Throttling is useful when replaying log segments to a downstream
// consumer that cannot handle burst traffic, or when simulating
// real-time log streaming at a controlled pace.
//
// A zero or negative delay disables the throttler, making all
// operations no-ops with no performance overhead.
//
// Example:
//
//	th := linethrottler.New(50 * time.Millisecond)
//	for _, line := range lines {
//		fmt.Println(th.Apply(line))
//	}
package linethrottler
