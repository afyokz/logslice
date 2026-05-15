// Package linecolumnextractor provides a positional column extractor for
// structured log lines.
//
// It splits each line by a configurable delimiter and maps the resulting
// fields to caller-supplied column names in declaration order. When the
// line contains more fields than there are column names the surplus text
// is collapsed into the last column, mirroring the behaviour of
// strings.SplitN.
//
// Usage:
//
//	e := linecolumnextractor.New(" ", []string{"timestamp", "level", "message"})
//	fields := e.Extract("2024-06-01T12:00:00Z INFO server started")
//	// fields["timestamp"] == "2024-06-01T12:00:00Z"
//	// fields["level"]     == "INFO"
//	// fields["message"]   == "server started"
package linecolumnextractor
