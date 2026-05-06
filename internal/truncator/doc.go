// Package truncator provides line-length limiting for log output.
//
// When exporting large log files, individual lines can be extremely long
// (e.g. JSON blobs, stack traces on a single line). The Truncator clips
// lines that exceed a configurable byte limit and optionally appends a
// visual marker such as "..." so consumers know the line was shortened.
//
// Usage:
//
//	tr := truncator.New(120, "...")
//	for _, line := range lines {
//		fmt.Fprintln(w, tr.Truncate(line))
//	}
//
// Truncation is disabled automatically when maxBytes is zero or negative,
// making it safe to wire into the pipeline unconditionally.
package truncator
