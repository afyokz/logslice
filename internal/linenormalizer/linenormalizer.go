// Package linenormalizer provides whitespace normalization for log lines.
// It collapses repeated whitespace characters into a single space and
// optionally trims leading and trailing whitespace.
package linenormalizer

import (
	"regexp"
	"strings"
)

var multiSpace = regexp.MustCompile(`\s+`)

// Normalizer collapses internal whitespace and optionally trims log lines.
type Normalizer struct {
	enabled bool
	trim    bool
}

// New creates a Normalizer. When enabled is false the normalizer is a no-op.
// When trim is true, leading and trailing whitespace is removed in addition
// to collapsing internal runs of whitespace.
func New(enabled, trim bool) *Normalizer {
	return &Normalizer{enabled: enabled, trim: trim}
}

// Enabled reports whether normalization is active.
func (n *Normalizer) Enabled() bool { return n.enabled }

// Trim reports whether edge trimming is active.
func (n *Normalizer) Trim() bool { return n.trim }

// Normalize returns the normalized form of line. If the normalizer is
// disabled the original string is returned unchanged.
func (n *Normalizer) Normalize(line string) string {
	if !n.enabled {
		return line
	}
	if n.trim {
		line = strings.TrimSpace(line)
	}
	return multiSpace.ReplaceAllString(line, " ")
}
