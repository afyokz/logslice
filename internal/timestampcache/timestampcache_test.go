package timestampcache_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/timestampcache"
)

func TestNew_DefaultCapacity(t *testing.T) {
	c := timestampcache.New(0)
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
}

func TestGet_MissOnEmpty(t *testing.T) {
	c := timestampcache.New(8)
	_, _, found := c.Get("2024-01-01T00:00:00Z")
	if found {
		t.Error("expected miss on empty cache")
	}
}

func TestPutAndGet_HitReturnsValue(t *testing.T) {
	c := timestampcache.New(8)
	now := time.Now().UTC().Truncate(time.Second)
	c.Put("2024-01-01T12:00:00Z", now, true)

	got, ok, found := c.Get("2024-01-01T12:00:00Z")
	if !found {
		t.Fatal("expected cache hit")
	}
	if !ok {
		t.Error("expected ok=true")
	}
	if !got.Equal(now) {
		t.Errorf("got %v, want %v", got, now)
	}
}

func TestPut_ParseFailureStored(t *testing.T) {
	c := timestampcache.New(8)
	c.Put("not-a-time", time.Time{}, false)

	_, ok, found := c.Get("not-a-time")
	if !found {
		t.Fatal("expected cache hit for failed parse")
	}
	if ok {
		t.Error("expected ok=false for failed parse")
	}
}

func TestEviction_OldestEvictedAtCapacity(t *testing.T) {
	const cap = 4
	c := timestampcache.New(cap)
	ts := time.Now().UTC()

	for i := 0; i < cap+1; i++ {
		key := time.Duration(i).String()
		c.Put(key, ts.Add(time.Duration(i)), true)
	}

	if c.Len() > cap {
		t.Errorf("cache len %d exceeds capacity %d", c.Len(), cap)
	}

	// The very first key should have been evicted.
	_, _, found := c.Get("0s")
	if found {
		t.Error("expected oldest entry to be evicted")
	}
}

func TestStats_HitsAndMisses(t *testing.T) {
	c := timestampcache.New(8)
	ts := time.Now().UTC()

	c.Get("miss1")
	c.Get("miss2")
	c.Put("hit1", ts, true)
	c.Get("hit1")
	c.Get("hit1")

	hits, misses := c.Stats()
	if hits != 2 {
		t.Errorf("hits: got %d, want 2", hits)
	}
	if misses != 2 {
		t.Errorf("misses: got %d, want 2", misses)
	}
}

func TestPut_DuplicateKeyNotOverwritten(t *testing.T) {
	c := timestampcache.New(8)
	t1 := time.Now().UTC()
	t2 := t1.Add(time.Hour)

	c.Put("key", t1, true)
	c.Put("key", t2, true) // should be a no-op

	got, _, _ := c.Get("key")
	if !got.Equal(t1) {
		t.Errorf("expected original value %v, got %v", t1, got)
	}
}
