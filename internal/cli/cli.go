// Package cli provides command-line interface parsing and configuration
// for the logslice tool.
package cli

import (
	"flag"
	"fmt"
	"os"
	"time"
)

// Config holds the parsed CLI arguments for a logslice run.
type Config struct {
	InputFile  string
	OutputFile string
	From       time.Time
	To         time.Time
	Format     string
	Numbered   bool
	Verbose    bool
}

// Parse parses os.Args and returns a Config or an error.
func Parse() (*Config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)

	input := fs.String("input", "", "Path to the input log file (required)")
	output := fs.String("output", "", "Path to the output file (default: stdout)")
	from := fs.String("from", "", "Start of time range, e.g. 2024-01-15T08:00:00")
	to := fs.String("to", "", "End of time range, e.g. 2024-01-15T09:00:00")
	format := fs.String("format", "", "Custom timestamp layout (Go time format)")
	numbered := fs.Bool("numbered", false, "Prefix output lines with line numbers")
	verbose := fs.Bool("verbose", false, "Print progress information to stderr")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	if *input == "" {
		return nil, fmt.Errorf("--input is required")
	}
	if *from == "" || *to == "" {
		return nil, fmt.Errorf("--from and --to are required")
	}

	layout := time.RFC3339
	if *format != "" {
		layout = *format
	}

	parsedFrom, err := time.Parse(layout, *from)
	if err != nil {
		return nil, fmt.Errorf("invalid --from value %q: %w", *from, err)
	}
	parsedTo, err := time.Parse(layout, *to)
	if err != nil {
		return nil, fmt.Errorf("invalid --to value %q: %w", *to, err)
	}
	if !parsedTo.After(parsedFrom) {
		return nil, fmt.Errorf("--to must be after --from")
	}

	return &Config{
		InputFile:  *input,
		OutputFile: *output,
		From:       parsedFrom,
		To:         parsedTo,
		Format:     layout,
		Numbered:   *numbered,
		Verbose:    *verbose,
	}, nil
}
