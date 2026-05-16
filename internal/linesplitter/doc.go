// Package linesplitter provides a simple field splitter that divides a log
// line into named parts using a configurable delimiter.
//
// It is designed for structured log formats where fields are separated by a
// known character or string (e.g. space, tab, pipe, comma) and the position
// of each field is fixed and known ahead of time.
//
// When fewer fields are configured than the delimiter produces the extra text
// is folded into the last field via strings.SplitN, preserving the remainder.
// When the line contains fewer parts than the number of configured field names
// the missing fields are returned as empty strings.
//
// Usage:
//
//	s := linesplitter.New(" ", []string{"level", "timestamp", "message"})
//	fields := s.Extract("INFO 2024-06-01T12:00:00Z application started")
//	// fields["level"]     == "INFO"
//	// fields["timestamp"] == "2024-06-01T12:00:00Z"
//	// fields["message"]   == "application started"
package linesplitter
