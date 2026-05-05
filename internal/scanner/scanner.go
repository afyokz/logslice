// Package scanner provides line-by-line log file scanning with time-range filtering.
package scanner

import (
	"bufio"
	"io"
	"time"

	"github.com/yourorg/logslice/internal/timeparser"
)

// Options configures the Scanner behavior.
type Options struct {
	// Format is the timestamp format string used to parse log lines.
	// If empty, ParseTimestamp auto-detection is used.
	Format string

	// Start is the beginning of the desired time range (inclusive).
	Start time.Time

	// End is the end of the desired time range (inclusive).
	End time.Time
}

// Scanner reads log lines from a reader and emits those within a time range.
type Scanner struct {
	r    io.Reader
	opts Options
}

// New creates a new Scanner that reads from r using the given Options.
func New(r io.Reader, opts Options) *Scanner {
	return &Scanner{r: r, opts: opts}
}

// Scan reads all lines from the underlying reader and writes matching lines to w.
// A line matches if its leading timestamp falls within [opts.Start, opts.End].
// Lines whose timestamp cannot be parsed are skipped.
// Returns the number of lines written and any read/write error.
func (s *Scanner) Scan(w io.Writer) (int, error) {
	br := bufio.NewReader(s.r)
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	written := 0
	for {
		line, err := br.ReadString('\n')
		if len(line) > 0 {
			var t time.Time
			var parseErr error

			if s.opts.Format != "" {
				t, parseErr = timeparser.ParseWithFormat(line, s.opts.Format)
			} else {
				t, parseErr = timeparser.ParseTimestamp(line)
			}

			if parseErr == nil && timeparser.InRange(t, s.opts.Start, s.opts.End) {
				if _, werr := bw.WriteString(line); werr != nil {
					return written, werr
				}
				written++
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return written, err
		}
	}

	return written, nil
}
