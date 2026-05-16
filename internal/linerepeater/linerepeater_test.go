package linerepeater

import (
	"testing"
)

func TestNew_DisabledWhenFalse(t *testing.T) {
	r := New(false, false)
	if r.Enabled() {
		t.Fatal("expected disabled")
	}
}

func TestNew_EnabledWhenTrue(t *testing.T) {
	r := New(true, false)
	if !r.Enabled() {
		t.Fatal("expected enabled")
	}
}

func TestFeed_Disabled_PassesThrough(t *testing.T) {
	r := New(false, false)
	for _, line := range []string{"a", "a", "b"} {
		out := r.Feed(line)
		if len(out) != 1 || out[0] != line {
			t.Fatalf("disabled: expected %q, got %v", line, out)
		}
	}
}

func TestFeed_Enabled_SuppressDuplicates(t *testing.T) {
	r := New(true, false)

	out := r.Feed("hello")
	if len(out) != 1 || out[0] != "hello" {
		t.Fatalf("first line: got %v", out)
	}

	// Duplicate lines should be suppressed.
	for i := 0; i < 3; i++ {
		out = r.Feed("hello")
		if len(out) != 0 {
			t.Fatalf("duplicate %d: expected empty, got %v", i, out)
		}
	}

	// New distinct line should appear.
	out = r.Feed("world")
	if len(out) != 1 || out[0] != "world" {
		t.Fatalf("new line: got %v", out)
	}
}

func TestFeed_Summary_EmitsCount(t *testing.T) {
	r := New(true, true)

	r.Feed("line")
	r.Feed("line") // duplicate
	r.Feed("line") // duplicate

	// Feeding a different line should trigger the summary.
	out := r.Feed("other")
	if len(out) != 2 {
		t.Fatalf("expected 2 outputs (summary + new line), got %v", out)
	}
	if out[0] != "... repeated 3 times" {
		t.Errorf("summary line: got %q", out[0])
	}
	if out[1] != "other" {
		t.Errorf("new line: got %q", out[1])
	}
}

func TestFlush_ReturnsSummary(t *testing.T) {
	r := New(true, true)

	r.Feed("x")
	r.Feed("x")
	r.Feed("x")

	out := r.Flush()
	if len(out) != 1 || out[0] != "... repeated 3 times" {
		t.Fatalf("flush: got %v", out)
	}
}

func TestFlush_NoSummary_WhenCountOne(t *testing.T) {
	r := New(true, true)
	r.Feed("only-once")
	out := r.Flush()
	if len(out) != 0 {
		t.Fatalf("expected no summary for single occurrence, got %v", out)
	}
}

func TestReset_ClearsState(t *testing.T) {
	r := New(true, true)
	r.Feed("a")
	r.Feed("a")
	r.Reset()

	// After reset, "a" should be treated as a fresh first occurrence.
	out := r.Feed("a")
	if len(out) != 1 || out[0] != "a" {
		t.Fatalf("after reset: got %v", out)
	}
	if r.Flush() != nil {
		t.Fatal("expected no summary after single occurrence post-reset")
	}
}
