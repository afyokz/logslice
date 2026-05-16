// Package linecounter provides a counter that tracks lines by category label.
package linecounter

// Counter accumulates line counts keyed by a string label.
// It is useful for producing per-category summaries at the end of a pipeline run.
type Counter struct {
	counts map[string]int
	enabled bool
}

// New creates a Counter. When enabled is false the counter is a no-op.
func New(enabled bool) *Counter {
	return &Counter{
		counts:  make(map[string]int),
		enabled: enabled,
	}
}

// Enabled reports whether the counter is active.
func (c *Counter) Enabled() bool { return c.enabled }

// Inc increments the count for label by 1.
// It is a no-op when the counter is disabled.
func (c *Counter) Inc(label string) {
	if !c.enabled {
		return
	}
	c.counts[label]++
}

// Add adds n to the count for label.
// It is a no-op when the counter is disabled or n <= 0.
func (c *Counter) Add(label string, n int) {
	if !c.enabled || n <= 0 {
		return
	}
	c.counts[label] += n
}

// Get returns the current count for label.
// It returns 0 for unknown labels or when disabled.
func (c *Counter) Get(label string) int {
	return c.counts[label]
}

// Labels returns all labels that have been recorded, in insertion order.
func (c *Counter) Labels() []string {
	labels := make([]string, 0, len(c.counts))
	for k := range c.counts {
		labels = append(labels, k)
	}
	return labels
}

// Reset clears all counts but preserves the enabled state.
func (c *Counter) Reset() {
	c.counts = make(map[string]int)
}

// Total returns the sum of all recorded counts.
func (c *Counter) Total() int {
	var t int
	for _, v := range c.counts {
		t += v
	}
	return t
}
