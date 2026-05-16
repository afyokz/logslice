package linechunker

import (
	"testing"
)

func TestNew_DisabledWhenZero(t *testing.T) {
	c := New(0)
	if c.Enabled() {
		t.Fatal("expected disabled for chunkSize=0")
	}
}

func TestNew_DisabledWhenNegative(t *testing.T) {
	c := New(-5)
	if c.Enabled() {
		t.Fatal("expected disabled for chunkSize=-5")
	}
}

func TestNew_EnabledWhenPositive(t *testing.T) {
	c := New(3)
	if !c.Enabled() {
		t.Fatal("expected enabled for chunkSize=3")
	}
	if c.ChunkSize() != 3 {
		t.Fatalf("ChunkSize: got %d, want 3", c.ChunkSize())
	}
}

func TestFeed_Disabled_PassesThrough(t *testing.T) {
	c := New(0)
	got := c.Feed("hello")
	if len(got) != 1 || got[0] != "hello" {
		t.Fatalf("disabled Feed: got %v, want [hello]", got)
	}
}

func TestFeed_AccumulatesUntilChunkSize(t *testing.T) {
	c := New(3)
	if got := c.Feed("a"); got != nil {
		t.Fatalf("expected nil after 1st line, got %v", got)
	}
	if got := c.Feed("b"); got != nil {
		t.Fatalf("expected nil after 2nd line, got %v", got)
	}
	got := c.Feed("c")
	if len(got) != 3 {
		t.Fatalf("expected chunk of 3, got %v", got)
	}
	if got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Fatalf("unexpected chunk contents: %v", got)
	}
}

func TestFeed_BufferResetAfterChunk(t *testing.T) {
	c := New(2)
	c.Feed("x")
	c.Feed("y") // flushes
	if got := c.Feed("z"); got != nil {
		t.Fatalf("expected nil for first line of new chunk, got %v", got)
	}
	got := c.Feed("w")
	if len(got) != 2 || got[0] != "z" || got[1] != "w" {
		t.Fatalf("second chunk wrong: %v", got)
	}
}

func TestFlush_ReturnsRemainder(t *testing.T) {
	c := New(5)
	c.Feed("one")
	c.Feed("two")
	got := c.Flush()
	if len(got) != 2 {
		t.Fatalf("Flush: got %v, want [one two]", got)
	}
}

func TestFlush_EmptyBuffer_ReturnsNil(t *testing.T) {
	c := New(3)
	if got := c.Flush(); got != nil {
		t.Fatalf("Flush on empty: got %v, want nil", got)
	}
}

func TestReset_DiscardsBuffer(t *testing.T) {
	c := New(4)
	c.Feed("a")
	c.Feed("b")
	c.Reset()
	if got := c.Flush(); got != nil {
		t.Fatalf("after Reset Flush should return nil, got %v", got)
	}
}

func TestFeed_MultipleFullChunks(t *testing.T) {
	c := New(2)
	var chunks [][]string
	for _, line := range []string{"1", "2", "3", "4", "5", "6"} {
		if ch := c.Feed(line); ch != nil {
			chunks = append(chunks, ch)
		}
	}
	if len(chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d: %v", len(chunks), chunks)
	}
}
