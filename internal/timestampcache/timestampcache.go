// Package timestampcache provides a fixed-size LRU cache for parsed log
// timestamps, avoiding redundant parsing of repeated or near-duplicate
// timestamp strings during high-volume log slicing.
package timestampcache

import (
	"sync"
	"time"
)

// entry holds a cached parse result.
type entry struct {
	t   time.Time
	ok  bool
	key string
}

// Cache is a thread-safe, fixed-capacity LRU cache mapping raw timestamp
// strings to parsed time.Time values.
type Cache struct {
	mu       sync.Mutex
	cap      int
	keys     []string // insertion-ordered ring for eviction
	pos      int      // next write position in ring
	store    map[string]entry
	hits     int
	misses   int
}

// New creates a Cache with the given capacity. If cap is less than 1 it
// defaults to 256.
func New(cap int) *Cache {
	if cap < 1 {
		cap = 256
	}
	return &Cache{
		cap:   cap,
		keys:  make([]string, cap),
		store: make(map[string]entry, cap),
	}
}

// Get returns the cached time.Time for raw and whether the entry was found
// and successfully parsed.
func (c *Cache) Get(raw string) (time.Time, bool, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, found := c.store[raw]; found {
		c.hits++
		return e.t, e.ok, true
	}
	c.misses++
	return time.Time{}, false, false
}

// Put stores a parse result for raw. If the cache is at capacity the oldest
// entry is evicted.
func (c *Cache) Put(raw string, t time.Time, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.store[raw]; exists {
		return
	}
	// Evict the entry that currently occupies the ring slot.
	if old := c.keys[c.pos]; old != "" {
		delete(c.store, old)
	}
	c.keys[c.pos] = raw
	c.pos = (c.pos + 1) % c.cap
	c.store[raw] = entry{t: t, ok: ok, key: raw}
}

// Stats returns the cumulative hit and miss counts.
func (c *Cache) Stats() (hits, misses int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.hits, c.misses
}

// Len returns the number of entries currently held in the cache.
func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.store)
}
