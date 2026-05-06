// Package cli handles command-line argument parsing for logslice.
//
// It exposes a single Parse function that reads os.Args, validates all
// required flags, and returns a Config struct ready for use by the
// scanner and exporter packages.
//
// Supported flags:
//
//	--input    path to the source log file (required)
//	--output   path for the output file; omit to write to stdout
//	--from     start of the time window in RFC3339 or custom format (required)
//	--to       end of the time window in RFC3339 or custom format (required)
//	--format   custom Go time layout used to parse --from / --to values
//	--numbered prefix each exported line with its original line number
//	--verbose  emit progress diagnostics to stderr
package cli
