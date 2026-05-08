// Package linefilter provides a simple length-based line filter that
// discards lines shorter or longer than configured thresholds.
package linefilter

// Filter discards lines that fall outside an optional min/max byte length.
type Filter struct {
	minLen int
	maxLen int
	enabled bool
}

// New creates a Filter. A minLen <= 0 means no lower bound; a maxLen <= 0
// means no upper bound. If neither bound is set the filter is disabled and
// Keep always returns true.
func New(minLen, maxLen int) *Filter {
	enabled := minLen > 0 || maxLen > 0
	return &Filter{
		minLen:  minLen,
		maxLen:  maxLen,
		enabled: enabled,
	}
}

// Enabled reports whether the filter has at least one active bound.
func (f *Filter) Enabled() bool { return f.enabled }

// Keep returns true when line satisfies the configured length constraints.
// If the filter is disabled Keep always returns true.
func (f *Filter) Keep(line string) bool {
	if !f.enabled {
		return true
	}
	n := len(line)
	if f.minLen > 0 && n < f.minLen {
		return false
	}
	if f.maxLen > 0 && n > f.maxLen {
		return false
	}
	return true
}

// Stats returns the min and max length thresholds.
func (f *Filter) Stats() (minLen, maxLen int) {
	return f.minLen, f.maxLen
}
