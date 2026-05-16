// Package lineskipper provides a configurable line skipper that drops
// the first N lines from a stream, useful for skipping file headers or
// preamble sections before processing log content.
package lineskipper

// Skipper drops the first N lines fed to it, passing all subsequent
// lines through unchanged.
type Skipper struct {
	enabled bool
	skip    int
	seen    int
}

// New returns a Skipper configured to drop the first n lines.
// If n is zero or negative, the Skipper is disabled and all lines pass through.
func New(n int) *Skipper {
	return &Skipper{
		enabled: n > 0,
		skip:    n,
	}
}

// Enabled reports whether the Skipper is active.
func (s *Skipper) Enabled() bool { return s.enabled }

// Skip returns the configured number of lines to drop.
func (s *Skipper) Skip() int { return s.skip }

// Keep reports whether the given line should be kept.
// The first Skip() lines return false; all subsequent lines return true.
// When disabled, Keep always returns true.
func (s *Skipper) Keep(line string) bool {
	if !s.enabled {
		return true
	}
	if s.seen < s.skip {
		s.seen++
		return false
	}
	return true
}

// Reset resets the internal counter so the next line is treated as the first.
func (s *Skipper) Reset() {
	s.seen = 0
}

// Count returns the number of lines that have been dropped so far.
func (s *Skipper) Count() int { return s.seen }
