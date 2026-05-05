// Package timeparser provides utilities for parsing and comparing timestamps
// found in log file lines.
//
// It supports a variety of common log timestamp formats including RFC3339,
// common syslog formats, and Apache/Nginx combined log formats. The package
// is used by logslice to identify whether a log line falls within a
// user-specified time range.
//
// Example usage:
//
//	t, format, err := timeparser.ParseTimestamp("2024-03-15T10:22:33Z")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("parsed with format:", format)
//
//	start, _ := time.Parse(time.RFC3339, "2024-03-15T00:00:00Z")
//	end, _ := time.Parse(time.RFC3339, "2024-03-15T23:59:59Z")
//	if timeparser.InRange(t, start, end) {
//		fmt.Println("line is within range")
//	}
package timeparser
