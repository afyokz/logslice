package truncator

// Truncator limits the byte length of log lines before export.
// Lines exceeding the maximum width are truncated and optionally
// suffixed with an ellipsis marker.
type Truncator struct {
	maxBytes int
	suffix   string
	enabled  bool
}

// New creates a Truncator with the given maximum byte width.
// A maxBytes value <= 0 disables truncation entirely.
// suffix is appended to truncated lines (e.g. "..."); pass "" for none.
func New(maxBytes int, suffix string) *Truncator {
	return &Truncator{
		maxBytes: maxBytes,
		suffix:   suffix,
		enabled:  maxBytes > 0,
	}
}

// Truncate returns the line unchanged if truncation is disabled or the
// line fits within the limit. Otherwise it returns a truncated copy.
func (t *Truncator) Truncate(line string) string {
	if !t.enabled {
		return line
	}
	if len(line) <= t.maxBytes {
		return line
	}
	cutAt := t.maxBytes
	if len(t.suffix) > 0 && len(t.suffix) < t.maxBytes {
		cutAt = t.maxBytes - len(t.suffix)
	}
	return line[:cutAt] + t.suffix
}

// Enabled reports whether truncation is active.
func (t *Truncator) Enabled() bool {
	return t.enabled
}

// MaxBytes returns the configured byte limit.
func (t *Truncator) MaxBytes() int {
	return t.maxBytes
}
