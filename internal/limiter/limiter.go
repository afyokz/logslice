package limiter

// Limiter enforces a maximum number of output lines.
// A zero or negative MaxLines value means no limit is applied.
type Limiter struct {
	MaxLines int
	count    int
}

// New creates a new Limiter with the given maximum line count.
// If maxLines is zero or negative, the limiter is effectively disabled.
func New(maxLines int) *Limiter {
	return &Limiter{MaxLines: maxLines}
}

// Allow reports whether the next line should be allowed through.
// It increments the internal counter on each accepted line.
// Once the limit is reached, all subsequent calls return false.
func (l *Limiter) Allow() bool {
	if l.MaxLines <= 0 {
		return true
	}
	if l.count >= l.MaxLines {
		return false
	}
	l.count++
	return true
}

// Count returns the number of lines that have been allowed so far.
func (l *Limiter) Count() int {
	return l.count
}

// Reset resets the internal counter, allowing the limiter to be reused.
func (l *Limiter) Reset() {
	l.count = 0
}

// Reached reports whether the limit has been reached.
func (l *Limiter) Reached() bool {
	if l.MaxLines <= 0 {
		return false
	}
	return l.count >= l.MaxLines
}
