// Package stats provides lightweight counters for tracking pipeline
// processing metrics such as lines read, matched, skipped, and exported.
package stats

import (
	"fmt"
	"io"
)

// Counter holds running totals accumulated during a pipeline run.
type Counter struct {
	LinesRead     int
	LinesMatched  int
	LinesSkipped  int
	LinesExported int
}

// Add merges another Counter into c.
func (c *Counter) Add(other Counter) {
	c.LinesRead += other.LinesRead
	c.LinesMatched += other.LinesMatched
	c.LinesSkipped += other.LinesSkipped
	c.LinesExported += other.LinesExported
}

// IncRead increments the read counter by 1.
func (c *Counter) IncRead() { c.LinesRead++ }

// IncMatched increments the matched counter by 1.
func (c *Counter) IncMatched() { c.LinesMatched++ }

// IncSkipped increments the skipped counter by 1.
func (c *Counter) IncSkipped() { c.LinesSkipped++ }

// IncExported increments the exported counter by 1.
func (c *Counter) IncExported() { c.LinesExported++ }

// Print writes a human-readable summary of the counter to w.
func (c *Counter) Print(w io.Writer) {
	fmt.Fprintf(w, "Lines read:     %d\n", c.LinesRead)
	fmt.Fprintf(w, "Lines matched:  %d\n", c.LinesMatched)
	fmt.Fprintf(w, "Lines skipped:  %d\n", c.LinesSkipped)
	fmt.Fprintf(w, "Lines exported: %d\n", c.LinesExported)
}
