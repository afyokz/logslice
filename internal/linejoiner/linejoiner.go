package linejoiner

import "strings"

// Joiner concatenates consecutive lines using a configurable delimiter,
// emitting a single joined line every N input lines.
type Joiner struct {
	enabled   bool
	every     int
	delimiter string
	buf       []string
}

// New returns a Joiner that groups every n lines and joins them with delim.
// If n <= 1 or delim is empty and n <= 1, the Joiner is disabled.
func New(n int, delim string) *Joiner {
	if n <= 1 {
		return &Joiner{enabled: false, every: 1, delimiter: delim}
	}
	if delim == "" {
		delim = " "
	}
	return &Joiner{
		enabled:   true,
		every:     n,
		delimiter: delim,
		buf:       make([]string, 0, n),
	}
}

// Enabled reports whether the joiner is active.
func (j *Joiner) Enabled() bool { return j.enabled }

// Feed accepts a line and returns a joined line when the group is complete,
// or an empty string and false if more lines are needed.
func (j *Joiner) Feed(line string) (string, bool) {
	if !j.enabled {
		return line, true
	}
	j.buf = append(j.buf, line)
	if len(j.buf) >= j.every {
		joined := strings.Join(j.buf, j.delimiter)
		j.buf = j.buf[:0]
		return joined, true
	}
	return "", false
}

// Flush returns any buffered lines joined together, or empty string if the
// buffer is empty. It resets the internal buffer.
func (j *Joiner) Flush() (string, bool) {
	if len(j.buf) == 0 {
		return "", false
	}
	joined := strings.Join(j.buf, j.delimiter)
	j.buf = j.buf[:0]
	return joined, true
}

// Reset clears the internal buffer without emitting anything.
func (j *Joiner) Reset() {
	j.buf = j.buf[:0]
}
