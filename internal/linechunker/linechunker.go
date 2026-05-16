package linechunker

// Chunker groups input lines into fixed-size batches (chunks). When disabled
// (chunkSize <= 0) every line is passed through as a single-element slice.
type Chunker struct {
	enabled   bool
	chunkSize int
	buf       []string
}

// New returns a Chunker that accumulates lines into batches of chunkSize.
// If chunkSize is zero or negative the Chunker is disabled.
func New(chunkSize int) *Chunker {
	return &Chunker{
		enabled:   chunkSize > 0,
		chunkSize: chunkSize,
	}
}

// Enabled reports whether chunking is active.
func (c *Chunker) Enabled() bool { return c.enabled }

// ChunkSize returns the configured batch size.
func (c *Chunker) ChunkSize() int { return c.chunkSize }

// Feed appends line to the internal buffer. When the buffer reaches chunkSize
// the accumulated lines are returned and the buffer is reset. Otherwise nil is
// returned, signalling that more lines are needed to complete the chunk.
//
// When disabled, Feed immediately returns a slice containing only line.
func (c *Chunker) Feed(line string) []string {
	if !c.enabled {
		return []string{line}
	}
	c.buf = append(c.buf, line)
	if len(c.buf) >= c.chunkSize {
		return c.flush()
	}
	return nil
}

// Flush returns whatever lines remain in the buffer and resets it.
// It returns nil when the buffer is empty.
func (c *Chunker) Flush() []string {
	if len(c.buf) == 0 {
		return nil
	}
	return c.flush()
}

// Reset discards all buffered lines without returning them.
func (c *Chunker) Reset() {
	c.buf = c.buf[:0]
}

func (c *Chunker) flush() []string {
	out := make([]string, len(c.buf))
	copy(out, c.buf)
	c.buf = c.buf[:0]
	return out
}
