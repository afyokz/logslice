package deduplicator

import "hash/fnv"

// Deduplicator filters out duplicate log lines within a sliding window.
type Deduplicator struct {
	seen     map[uint64]struct{}
	window   int
	queue    []uint64
	enabled  bool
}

// New creates a Deduplicator with the given window size.
// If window <= 0, deduplication is disabled.
func New(window int) *Deduplicator {
	if window <= 0 {
		return &Deduplicator{enabled: false}
	}
	return &Deduplicator{
		seen:   make(map[uint64]struct{}, window),
		window: window,
		queue:  make([]uint64, 0, window),
		enabled: true,
	}
}

// IsDuplicate returns true if the line was seen within the current window.
// It also records the line for future lookups.
func (d *Deduplicator) IsDuplicate(line string) bool {
	if !d.enabled {
		return false
	}
	h := hash(line)
	if _, exists := d.seen[h]; exists {
		return true
	}
	d.record(h)
	return false
}

// Reset clears all seen entries.
func (d *Deduplicator) Reset() {
	if !d.enabled {
		return
	}
	d.seen = make(map[uint64]struct{}, d.window)
	d.queue = d.queue[:0]
}

// Enabled reports whether deduplication is active.
func (d *Deduplicator) Enabled() bool {
	return d.enabled
}

func (d *Deduplicator) record(h uint64) {
	if len(d.queue) >= d.window {
		old := d.queue[0]
		d.queue = d.queue[1:]
		delete(d.seen, old)
	}
	d.seen[h] = struct{}{}
	d.queue = append(d.queue, h)
}

func hash(s string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
}
