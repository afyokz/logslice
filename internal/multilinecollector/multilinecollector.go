package multilinecollector

import (
	"strings"

	"github.com/yourorg/logslice/internal/headerdetector"
)

// Collector groups continuation lines with their header line.
// Lines that do not start a new log entry (e.g. stack-trace frames) are
// appended to the previous entry so the pipeline sees one logical record.
type Collector struct {
	detector  *headerdetector.Detector
	pending   []string
	separator string
}

// New returns a Collector that uses det to decide whether a line is a
// header. separator is placed between joined lines (typically "\n").
func New(det *headerdetector.Detector, separator string) *Collector {
	if separator == "" {
		separator = "\n"
	}
	return &Collector{detector: det, separator: separator}
}

// Feed accepts the next raw line. It returns a complete (flushed) log
// record and true when a new header is encountered and the previously
// accumulated lines are ready, or "", false when lines are still being
// collected.
func (c *Collector) Feed(line string) (string, bool) {
	if c.detector.IsHeader(line) {
		flushed, ok := c.flush()
		c.pending = []string{line}
		return flushed, ok
	}
	// continuation line — accumulate
	c.pending = append(c.pending, line)
	return "", false
}

// Flush returns whatever has been accumulated so far and resets state.
// Call this after the input is exhausted to retrieve the final record.
func (c *Collector) Flush() (string, bool) {
	return c.flush()
}

// Reset discards any accumulated state.
func (c *Collector) Reset() {
	c.pending = c.pending[:0]
}

func (c *Collector) flush() (string, bool) {
	if len(c.pending) == 0 {
		return "", false
	}
	result := strings.Join(c.pending, c.separator)
	c.pending = c.pending[:0]
	return result, true
}
