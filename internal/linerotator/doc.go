// Package linerotator provides a line-aware writer that automatically rotates
// output to a new io.Writer after a configurable number of lines have been
// written.
//
// This is useful when exporting large time-range slices that should be split
// into smaller segment files, e.g. one file per N log lines.
//
// Usage:
//
//	seq := 0
//	factory := func(n int) (io.Writer, error) {
//		return os.Create(fmt.Sprintf("segment-%03d.log", n))
//	}
//	r, err := linerotator.New(nil, factory, 1000)
//	if err != nil { ... }
//	for _, line := range lines {
//		r.WriteLine(line)
//	}
//
// When maxLines is zero or negative the Rotator is disabled and all output is
// forwarded to the fallback writer supplied to New.
package linerotator
