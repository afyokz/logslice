package linebuffer

// Buffer is a fixed-capacity ring buffer that retains the last N lines.
// It is useful for capturing trailing context around matched log entries.
type Buffer struct {
	lines    []string
	cap      int
	head     int
	size     int
	enabled  bool
}

// New creates a new Buffer with the given capacity.
// If capacity is <= 0, the buffer is disabled and all operations are no-ops.
func New(capacity int) *Buffer {
	if capacity <= 0 {
		return &Buffer{enabled: false}
	}
	return &Buffer{
		lines:   make([]string, capacity),
		cap:     capacity,
		enabled: true,
	}
}

// Enabled reports whether the buffer is active.
func (b *Buffer) Enabled() bool { return b.enabled }

// Push adds a line to the ring buffer, evicting the oldest entry when full.
func (b *Buffer) Push(line string) {
	if !b.enabled {
		return
	}
	b.lines[b.head] = line
	b.head = (b.head + 1) % b.cap
	if b.size < b.cap {
		b.size++
	}
}

// Lines returns all buffered lines in insertion order (oldest first).
// Returns nil when the buffer is disabled or empty.
func (b *Buffer) Lines() []string {
	if !b.enabled || b.size == 0 {
		return nil
	}
	out := make([]string, b.size)
	start := (b.head - b.size + b.cap) % b.cap
	for i := 0; i < b.size; i++ {
		out[i] = b.lines[(start+i)%b.cap]
	}
	return out
}

// Reset clears all buffered lines without releasing memory.
func (b *Buffer) Reset() {
	if !b.enabled {
		return
	}
	b.head = 0
	b.size = 0
}

// Len returns the number of lines currently held in the buffer.
func (b *Buffer) Len() int { return b.size }
