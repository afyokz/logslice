package multilinecollector_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/headerdetector"
	"github.com/yourorg/logslice/internal/multilinecollector"
)

func newCollector(t *testing.T) *multilinecollector.Collector {
	t.Helper()
	det, err := headerdetector.New(nil)
	if err != nil {
		t.Fatalf("headerdetector.New: %v", err)
	}
	return multilinecollector.New(det, "\n")
}

func TestFeed_HeaderLine_FlushesNothing_OnFirst(t *testing.T) {
	c := newCollector(t)
	_, ok := c.Feed("2024-01-01T00:00:00Z INFO starting")
	if ok {
		t.Error("expected no flush on first header")
	}
}

func TestFeed_SecondHeader_FlushesPrevious(t *testing.T) {
	c := newCollector(t)
	c.Feed("2024-01-01T00:00:00Z INFO first") //nolint
	got, ok := c.Feed("2024-01-01T00:00:01Z INFO second")
	if !ok {
		t.Fatal("expected flush on second header")
	}
	if got != "2024-01-01T00:00:00Z INFO first" {
		t.Errorf("unexpected flushed record: %q", got)
	}
}

func TestFeed_ContinuationLines_AppendedToHeader(t *testing.T) {
	c := newCollector(t)
	c.Feed("2024-01-01T00:00:00Z ERROR boom") //nolint
	c.Feed("\tat foo.Bar(foo.go:42)")          //nolint
	c.Feed("\tat baz.Qux(baz.go:7)")          //nolint
	got, ok := c.Feed("2024-01-01T00:00:01Z INFO recovered")
	if !ok {
		t.Fatal("expected flush")
	}
	want := "2024-01-01T00:00:00Z ERROR boom\n\tat foo.Bar(foo.go:42)\n\tat baz.Qux(baz.go:7)"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFlush_ReturnsLastRecord(t *testing.T) {
	c := newCollector(t)
	c.Feed("2024-01-01T00:00:00Z INFO only") //nolint
	got, ok := c.Flush()
	if !ok {
		t.Fatal("expected flush to return record")
	}
	if got != "2024-01-01T00:00:00Z INFO only" {
		t.Errorf("unexpected record: %q", got)
	}
}

func TestFlush_EmptyCollector_ReturnsFalse(t *testing.T) {
	c := newCollector(t)
	_, ok := c.Flush()
	if ok {
		t.Error("expected no record from empty collector")
	}
}

func TestReset_ClearsState(t *testing.T) {
	c := newCollector(t)
	c.Feed("2024-01-01T00:00:00Z INFO something") //nolint
	c.Reset()
	_, ok := c.Flush()
	if ok {
		t.Error("expected empty after reset")
	}
}

func TestNew_DefaultSeparator(t *testing.T) {
	det, _ := headerdetector.New(nil)
	c := multilinecollector.New(det, "")
	c.Feed("2024-01-01T00:00:00Z INFO a") //nolint
	c.Feed("\tcontinuation")              //nolint
	got, _ := c.Flush()
	if got != "2024-01-01T00:00:00Z INFO a\n\tcontinuation" {
		t.Errorf("unexpected: %q", got)
	}
}
