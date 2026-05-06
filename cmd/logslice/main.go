package main

import (
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/cli"
	"github.com/yourorg/logslice/internal/pipeline"
)

const version = "0.1.0"

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("logslice version %s\n", version)
		os.Exit(0)
	}

	cfg, err := cli.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		printUsage()
		os.Exit(1)
	}

	if err := pipeline.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `
Usage:
  logslice --input <file> --from <timestamp> --to <timestamp> [options]

Options:
  --input   Path to the log file to slice
  --from    Start of time range (e.g. "2024-01-15 08:00:00")
  --to      End of time range (e.g. "2024-01-15 09:00:00")
  --output  Output file path (default: stdout)
  --format  Timestamp format (default: auto-detect)
  --include Comma-separated substrings or patterns to include
  --exclude Comma-separated substrings or patterns to exclude
  --number  Prefix output lines with line numbers
  --version Print version and exit

Examples:
  logslice --input app.log --from "2024-01-15 08:00:00" --to "2024-01-15 09:00:00"
  logslice --input app.log --from "2024-01-15 08:00:00" --to "2024-01-15 09:00:00" --include ERROR,WARN --output out.log`)
}
