// Package deduplicator provides a sliding-window duplicate line filter
// for use in the logslice pipeline.
//
// Lines are hashed using FNV-64a and tracked within a fixed-size window.
// When the window is full, the oldest entry is evicted to make room for
// new ones. This bounds memory usage while still catching consecutive
// duplicate log entries that commonly appear in high-throughput logs.
//
// Usage:
//
//	d := deduplicator.New(1000) // track last 1000 unique lines
//	for _, line := range lines {
//		if !d.IsDuplicate(line) {
//			// emit line
//		}
//	}
//
// Setting window to 0 or a negative value disables deduplication entirely
// and IsDuplicate always returns false, incurring no overhead.
package deduplicator
