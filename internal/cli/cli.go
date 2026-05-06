// Package cli parses and validates command-line arguments for logslice.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// Args holds the parsed and validated CLI arguments.
type Args struct {
	Input           string
	Output          string
	From            time.Time
	To              time.Time
	Format          string
	Numbered        bool
	IncludePatterns []string
	ExcludePatterns []string
}

// Parse parses os.Args and returns validated Args or an error.
func Parse() (*Args, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	input := fs.String("input", "", "path to the log file (required)")
	output := fs.String("output", "", "output file path (default: stdout)")
	from := fs.String("from", "", "start timestamp, e.g. 2006-01-02T15:04:05 (required)")
	to := fs.String("to", "", "end timestamp, e.g. 2006-01-02T15:04:05 (required)")
	format := fs.String("format", time.RFC3339, "timestamp format used in the log file")
	numbered := fs.Bool("numbered", false, "prefix output lines with line numbers")
	include := fs.String("include", "", "comma-separated regex patterns to include")
	exclude := fs.String("exclude", "", "comma-separated regex patterns to exclude")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	if *input == "" {
		return nil, errors.New("--input is required")
	}
	if *from == "" {
		return nil, errors.New("--from is required")
	}
	if *to == "" {
		return nil, errors.New("--to is required")
	}

	fromTime, err := time.Parse(*format, *from)
	if err != nil {
		return nil, fmt.Errorf("invalid --from value %q: %w", *from, err)
	}

	toTime, err := time.Parse(*format, *to)
	if err != nil {
		return nil, fmt.Errorf("invalid --to value %q: %w", *to, err)
	}

	if !toTime.After(fromTime) {
		return nil, errors.New("--to must be after --from")
	}

	args := &Args{
		Input:    *input,
		Output:   *output,
		From:     fromTime,
		To:       toTime,
		Format:   *format,
		Numbered: *numbered,
	}

	if *include != "" {
		args.IncludePatterns = splitTrimmed(*include)
	}
	if *exclude != "" {
		args.ExcludePatterns = splitTrimmed(*exclude)
	}

	return args, nil
}

func splitTrimmed(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
