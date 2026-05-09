package linesorter_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/linesorter"
)

// mockParse extracts a unix-second timestamp encoded as the first token of the
// line in the form "t=<unix>" for test convenience.
func mockParse(line string) (time.Time, error) {
	var sec int64
	_, err := fmt.Sscanf(line, "t=%d", &sec)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(sec, 0), nil
}

func TestNew_DisabledWhenParseNil(t *testing.T) {
	s := linesorter.New(true, nil)
	if s.Enabled() {
		t.Fatal("expected disabled when parse func is nil")
	}
}

func TestNew_DisabledExplicitly(t *testing.T) {
	s := linesorter.New(false, mockParse)
	if s.Enabled() {
		t.Fatal("expected disabled")
	}
}

func TestNew_EnabledWhenParseProvided(t *testing.T) {
	s := linesorter.New(true, mockParse)
	if !s.Enabled() {
		t.Fatal("expected enabled")
	}
}

func TestFlush_Disabled_PreservesInsertionOrder(t *testing.T) {
	s := linesorter.New(false, mockParse)
	lines := []string{"t=3 c", "t=1 a", "t=2 b"}
	for _, l := range lines {
		s.Feed(l)
	}
	got := s.Flush()
	for i, want := range lines {
		if got[i] != want {
			t.Fatalf("index %d: got %q, want %q", i, got[i], want)
		}
	}
}

func TestFlush_Enabled_SortsByTimestamp(t *testing.T) {
	s := linesorter.New(true, mockParse)
	s.Feed("t=30 third")
	s.Feed("t=10 first")
	s.Feed("t=20 second")

	got := s.Flush()
	want := []string{"t=10 first", "t=20 second", "t=30 third"}
	for i, w := range want {
		if got[i] != w {
			t.Fatalf("index %d: got %q, want %q", i, got[i], w)
		}
	}
}

func TestFlush_StableForEqualTimestamps(t *testing.T) {
	s := linesorter.New(true, mockParse)
	s.Feed("t=5 alpha")
	s.Feed("t=5 beta")
	s.Feed("t=5 gamma")

	got := s.Flush()
	want := []string{"t=5 alpha", "t=5 beta", "t=5 gamma"}
	for i, w := range want {
		if got[i] != w {
			t.Fatalf("index %d: got %q, want %q", i, got[i], w)
		}
	}
}

func TestFlush_ResetsBuffer(t *testing.T) {
	s := linesorter.New(true, mockParse)
	s.Feed("t=1 x")
	s.Flush()
	if s.Len() != 0 {
		t.Fatalf("expected empty buffer after Flush, got %d", s.Len())
	}
}

func TestFlush_EmptyBuffer_ReturnsEmptySlice(t *testing.T) {
	s := linesorter.New(true, mockParse)
	got := s.Flush()
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestLen_TracksBufferedLines(t *testing.T) {
	s := linesorter.New(true, mockParse)
	if s.Len() != 0 {
		t.Fatal("expected 0 initially")
	}
	s.Feed("t=1 a")
	s.Feed("t=2 b")
	if s.Len() != 2 {
		t.Fatalf("expected 2, got %d", s.Len())
	}
}
