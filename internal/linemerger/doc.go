// Package linemerger provides a Merger that groups consecutive log lines
// sharing a common leading prefix into a single logical record.
//
// # Overview
//
// Some log formats emit multi-part entries where each physical line starts
// with the same timestamp prefix. Merger detects prefix boundaries and joins
// the constituent lines with a configurable delimiter so that downstream
// consumers receive one logical record per group.
//
// # Usage
//
//	m := linemerger.New(19, " | ") // group by first 19 chars (timestamp)
//	for _, line := range lines {
//		if out, ok := m.Feed(line); ok {
//			process(out)
//		}
//	}
//	if out, ok := m.Flush(); ok {
//		process(out)
//	}
package linemerger
