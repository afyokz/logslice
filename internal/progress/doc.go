// Package progress provides a lightweight progress reporter for tracking
// log file processing metrics during a logslice run.
//
// A Reporter tracks three counters:
//   - Processed: total lines read from the input
//   - Matched:   lines that fell within the requested time range and
//                passed any include/exclude filters
//   - Skipped:   lines whose timestamps could not be parsed
//
// When verbose mode is enabled and a total line count is provided,
// Print also emits a percentage-complete indicator.
//
// Example usage:
//
//	rep := progress.New(totalLines, verbose, os.Stderr)
//	// ... for each line processed:
//	rep.IncProcessed()
//	rep.IncMatched() // or IncSkipped()
//	rep.Print()
package progress
