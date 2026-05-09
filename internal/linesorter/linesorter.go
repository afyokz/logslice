// Package linesorter provides a buffer-based sorter that collects log lines
// and emits them ordered by their parsed timestamp.
package linesorter

import (
	"sort"
	"time"
)

// entry holds a single log line paired with its parsed timestamp.
type entry struct {
	line string
	ts   time.Time
}

// Sorter accumulates lines and can flush them in timestamp order.
type Sorter struct {
	entries  []entry
	parse    func(string) (time.Time, error)
	enabled  bool
}

// New returns a Sorter. When parse is nil or enabled is false the Sorter is a
// no-op pass-through: Feed appends lines as-is and Flush returns them in
// insertion order.
func New(enabled bool, parse func(string) (time.Time, error)) *Sorter {
	if parse == nil {
		enabled = false
	}
	return &Sorter{
		parse:   parse,
		enabled: enabled,
	}
}

// Enabled reports whether timestamp-based sorting is active.
func (s *Sorter) Enabled() bool { return s.enabled }

// Feed adds a line to the internal buffer.
// If the sorter is disabled the line is stored with a zero timestamp.
func (s *Sorter) Feed(line string) {
	var ts time.Time
	if s.enabled {
		if t, err := s.parse(line); err == nil {
			ts = t
		}
	}
	s.entries = append(s.entries, entry{line: line, ts: ts})
}

// Flush sorts the buffered lines by timestamp (stable, preserving insertion
// order for equal timestamps) and returns them. The internal buffer is reset.
func (s *Sorter) Flush() []string {
	if s.enabled {
		sort.SliceStable(s.entries, func(i, j int) bool {
			return s.entries[i].ts.Before(s.entries[j].ts)
		})
	}
	out := make([]string, len(s.entries))
	for i, e := range s.entries {
		out[i] = e.line
	}
	s.entries = s.entries[:0]
	return out
}

// Len returns the number of lines currently buffered.
func (s *Sorter) Len() int { return len(s.entries) }
