// Package fieldextractor provides a configurable field extraction utility
// for structured log lines.
//
// It uses named regular expression capture groups to parse key fields
// (such as log level, timestamp, message, or IP address) from raw log lines.
//
// Usage:
//
//	e, err := fieldextractor.New(`(?P<level>\w+)\s+(?P<msg>.+)`)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fields := e.Extract("ERROR something went wrong")
//	// fields["level"] == "ERROR"
//	// fields["msg"]   == "something went wrong"
//
// If the pattern is empty, the extractor is disabled and Extract returns nil.
package fieldextractor
