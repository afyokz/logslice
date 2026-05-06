// Package exporter handles writing filtered log lines to output destinations.
package exporter

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Format represents the output format for exported log lines.
type Format int

const (
	// FormatRaw writes lines as-is.
	FormatRaw Format = iota
	// FormatNumbered prefixes each line with its original line number.
	FormatNumbered
)

// Options configures the exporter behavior.
type Options struct {
	Format     Format
	OutputPath string // empty means stdout
}

// Exporter writes log lines to a configured destination.
type Exporter struct {
	opts   Options
	writer io.Writer
	closer io.Closer
}

// New creates a new Exporter. Call Close when done.
func New(opts Options) (*Exporter, error) {
	if opts.OutputPath == "" {
		return &Exporter{opts: opts, writer: os.Stdout}, nil
	}

	f, err := os.Create(opts.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("exporter: create output file: %w", err)
	}
	return &Exporter{opts: opts, writer: f, closer: f}, nil
}

// Export writes the provided lines to the configured output.
func (e *Exporter) Export(lines []string) error {
	bw := bufio.NewWriter(e.writer)
	for i, line := range lines {
		var out string
		switch e.opts.Format {
		case FormatNumbered:
			out = fmt.Sprintf("%d: %s\n", i+1, line)
		default:
			out = line + "\n"
		}
		if _, err := bw.WriteString(out); err != nil {
			return fmt.Errorf("exporter: write line: %w", err)
		}
	}
	return bw.Flush()
}

// Close releases any resources held by the exporter.
func (e *Exporter) Close() error {
	if e.closer != nil {
		return e.closer.Close()
	}
	return nil
}
