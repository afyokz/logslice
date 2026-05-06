package progress

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"
)

// Reporter tracks and reports progress of log file processing.
type Reporter struct {
	out        io.Writer
	total      int64
	processed  atomic.Int64
	matched    atomic.Int64
	skipped    atomic.Int64
	verbose    bool
}

// New creates a new Reporter. If out is nil, os.Stderr is used.
func New(total int64, verbose bool, out io.Writer) *Reporter {
	if out == nil {
		out = os.Stderr
	}
	return &Reporter{
		out:     out,
		total:   total,
		verbose: verbose,
	}
}

// IncProcessed increments the count of processed lines.
func (r *Reporter) IncProcessed() {
	r.processed.Add(1)
}

// IncMatched increments the count of matched (exported) lines.
func (r *Reporter) IncMatched() {
	r.matched.Add(1)
}

// IncSkipped increments the count of skipped (unparsable) lines.
func (r *Reporter) IncSkipped() {
	r.skipped.Add(1)
}

// Print writes a summary of the current progress to the output writer.
func (r *Reporter) Print() {
	processed := r.processed.Load()
	matched := r.matched.Load()
	skipped := r.skipped.Load()

	if r.verbose {
		fmt.Fprintf(r.out, "progress: processed=%d matched=%d skipped=%d",
			processed, matched, skipped)
		if r.total > 0 {
			pct := float64(processed) / float64(r.total) * 100
			fmt.Fprintf(r.out, " total=%d (%.1f%%)", r.total, pct)
		}
		fmt.Fprintln(r.out)
	} else {
		fmt.Fprintf(r.out, "processed=%d matched=%d skipped=%d\n",
			processed, matched, skipped)
	}
}

// Summary returns a formatted summary string.
func (r *Reporter) Summary() string {
	return fmt.Sprintf("processed=%d matched=%d skipped=%d",
		r.processed.Load(), r.matched.Load(), r.skipped.Load())
}
