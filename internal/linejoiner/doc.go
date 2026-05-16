// Package linejoiner provides a Joiner that groups consecutive log lines and
// concatenates them into a single line using a configurable delimiter.
//
// This is useful when log entries are split across a fixed number of physical
// lines and must be reassembled before further processing, such as timestamp
// parsing or pattern matching.
//
// Basic usage:
//
//	j := linejoiner.New(3, " | ")
//	for _, raw := range inputLines {
//		if joined, ok := j.Feed(raw); ok {
//			process(joined)
//		}
//	}
//	if remainder, ok := j.Flush(); ok {
//		process(remainder)
//	}
//
// When n <= 1 the Joiner is disabled and Feed passes every line through
// unchanged, making it safe to include unconditionally in a pipeline.
package linejoiner
