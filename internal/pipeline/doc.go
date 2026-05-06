// Package pipeline wires together the scanner, matcher, and exporter
// into a single cohesive execution unit.
//
// A [Config] value describes the full set of options for a slice run:
// the input reader, output writer, time boundaries, include/exclude
// patterns, and export formatting preferences.
//
// [Run] executes the three-stage pipeline:
//
//  1. Scan — reads the input and retains only lines whose timestamps
//     fall within [From, To].
//  2. Match — applies include/exclude regex patterns to the scanned lines.
//  3. Export — writes the surviving lines to the configured destination,
//     optionally numbering them.
//
// The returned [Result] reports how many lines were scanned, matched,
// and ultimately exported, which is useful for progress reporting and
// testing.
package pipeline
