package timeparser

import (
	"fmt"
	"time"
)

// Common log timestamp formats to attempt parsing
var knownFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04:05.000",
	"2006-01-02 15:04:05.999999999",
	"02/Jan/2006:15:04:05 -0700",
	"Jan 2 15:04:05",
	"Jan 02 15:04:05",
}

// ParseTimestamp attempts to parse a timestamp string using known formats.
// Returns the parsed time and the matched format string, or an error if none match.
func ParseTimestamp(s string) (time.Time, string, error) {
	for _, format := range knownFormats {
		t, err := time.Parse(format, s)
		if err == nil {
			return t, format, nil
		}
	}
	return time.Time{}, "", fmt.Errorf("timeparser: unable to parse timestamp %q", s)
}

// ParseWithFormat parses a timestamp using a specific format.
func ParseWithFormat(s, format string) (time.Time, error) {
	t, err := time.Parse(format, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("timeparser: parse with format %q failed: %w", format, err)
	}
	return t, nil
}

// InRange reports whether t falls within [start, end] (inclusive).
func InRange(t, start, end time.Time) bool {
	return !t.Before(start) && !t.After(end)
}
