package linethrottler

import (
	"time"
)

// Throttler introduces a configurable delay between emitted lines,
// allowing rate-controlled output for downstream consumers.
type Throttler struct {
	enabled  bool
	delay    time.Duration
	sleepFn  func(time.Duration)
}

// New returns a Throttler that pauses for delay before each line.
// A zero or negative delay disables throttling.
func New(delay time.Duration) *Throttler {
	t := &Throttler{
		enabled: delay > 0,
		delay:   delay,
		sleepFn: time.Sleep,
	}
	return t
}

// Enabled reports whether throttling is active.
func (t *Throttler) Enabled() bool {
	return t.enabled
}

// Delay returns the configured inter-line delay.
func (t *Throttler) Delay() time.Duration {
	return t.delay
}

// Wait blocks for the configured delay if throttling is enabled.
// It is a no-op when disabled.
func (t *Throttler) Wait() {
	if !t.enabled {
		return
	}
	t.sleepFn(t.delay)
}

// Apply calls Wait and then returns the line unchanged.
// It is a convenience wrapper for pipeline use.
func (t *Throttler) Apply(line string) string {
	t.Wait()
	return line
}
