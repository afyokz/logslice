package linewindow

// Window holds a fixed-size sliding buffer of the most recently seen lines.
// It is used to emit N lines of context before or after a matching log entry.
type Window struct {
	buf  []string
	size int
	head int
	count int
}

// New returns a Window that retains up to size lines.
// If size is <= 0 the window is disabled (zero capacity).
func New(size int) *Window {
	if size <= 0 {
		size = 0
	}
	buf := make([]string, size)
	return &Window{buf: buf, size: size}
}

// Enabled reports whether the window has a non-zero capacity.
func (w *Window) Enabled() bool { return w.size > 0 }

// Push adds a line to the window, evicting the oldest entry when full.
func (w *Window) Push(line string) {
	if w.size == 0 {
		return
	}
	w.buf[w.head] = line
	w.head = (w.head + 1) % w.size
	if w.count < w.size {
		w.count++
	}
}

// Lines returns the buffered lines in insertion order (oldest first).
// The returned slice is a copy.
func (w *Window) Lines() []string {
	if w.count == 0 {
		return nil
	}
	out := make([]string, w.count)
	start := (w.head - w.count + w.size) % w.size
	for i := 0; i < w.count; i++ {
		out[i] = w.buf[(start+i)%w.size]
	}
	return out
}

// Reset clears all buffered lines without reallocating.
func (w *Window) Reset() {
	w.head = 0
	w.count = 0
}

// Len returns the number of lines currently buffered.
func (w *Window) Len() int { return w.count }
