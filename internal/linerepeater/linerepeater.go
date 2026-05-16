// Package linerepeater suppresses repeated consecutive identical lines,
// optionally emitting a summary count when a run of duplicates ends.
package linerepeater

import "fmt"

// Repeater tracks consecutive identical lines and can suppress or
// summarise runs of duplicates.
type Repeater struct {
	enabled bool
	summary bool
	last    string
	count   int
}

// New creates a Repeater. When enabled is false the Repeater is a no-op.
// When summary is true a "... repeated N times" line is emitted after a run.
func New(enabled, summary bool) *Repeater {
	return &Repeater{enabled: enabled, summary: summary}
}

// Enabled reports whether the repeater is active.
func (r *Repeater) Enabled() bool { return r.enabled }

// Feed accepts the next log line and returns the lines that should be
// emitted at this point (zero, one, or two entries).
//
// When disabled every line is returned as-is.
func (r *Repeater) Feed(line string) []string {
	if !r.enabled {
		return []string{line}
	}

	if line == r.last {
		r.count++
		return nil
	}

	// Line differs from the previous one – flush any pending summary first.
	out := r.flush()

	r.last = line
	r.count = 1
	out = append(out, line)
	return out
}

// Flush finalises any in-progress run and returns pending summary lines.
// Call after the last line has been fed.
func (r *Repeater) Flush() []string {
	if !r.enabled {
		return nil
	}
	return r.flush()
}

// flush emits a summary if the previous run had more than one occurrence.
func (r *Repeater) flush() []string {
	if r.summary && r.count > 1 {
		return []string{fmt.Sprintf("... repeated %d times", r.count)}
	}
	return nil
}

// Reset clears internal state so the Repeater can be reused.
func (r *Repeater) Reset() {
	r.last = ""
	r.count = 0
}
