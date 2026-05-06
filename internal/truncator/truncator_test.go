package truncator

import (
	"strings"
	"testing"
)

func TestNew_DisabledWhenZero(t *testing.T) {
	tr := New(0, "...")
	if tr.Enabled() {
		t.Fatal("expected truncator to be disabled for maxBytes=0")
	}
}

func TestNew_DisabledWhenNegative(t *testing.T) {
	tr := New(-5, "...")
	if tr.Enabled() {
		t.Fatal("expected truncator to be disabled for negative maxBytes")
	}
}

func TestNew_EnabledWhenPositive(t *testing.T) {
	tr := New(80, "...")
	if !tr.Enabled() {
		t.Fatal("expected truncator to be enabled for maxBytes=80")
	}
	if tr.MaxBytes() != 80 {
		t.Fatalf("expected MaxBytes=80, got %d", tr.MaxBytes())
	}
}

func TestTruncate_Disabled_ReturnsOriginal(t *testing.T) {
	tr := New(0, "...")
	line := strings.Repeat("x", 200)
	if got := tr.Truncate(line); got != line {
		t.Fatalf("expected original line, got truncated")
	}
}

func TestTruncate_ShortLine_Unchanged(t *testing.T) {
	tr := New(80, "...")
	line := "short line"
	if got := tr.Truncate(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestTruncate_ExactLength_Unchanged(t *testing.T) {
	tr := New(10, "")
	line := "1234567890"
	if got := tr.Truncate(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestTruncate_LongLine_WithSuffix(t *testing.T) {
	tr := New(10, "...")
	line := "abcdefghijklmnopqrstuvwxyz"
	got := tr.Truncate(line)
	if len(got) > 10 {
		t.Fatalf("expected len <= 10, got %d: %q", len(got), got)
	}
	if !strings.HasSuffix(got, "...") {
		t.Fatalf("expected suffix '...', got %q", got)
	}
	if got != "abcdefg..." {
		t.Fatalf("expected %q, got %q", "abcdefg...", got)
	}
}

func TestTruncate_LongLine_NoSuffix(t *testing.T) {
	tr := New(5, "")
	line := "hello world"
	got := tr.Truncate(line)
	if got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestTruncate_SuffixLargerThanLimit_CutsAtMax(t *testing.T) {
	tr := New(3, "...")
	line := "abcdefgh"
	got := tr.Truncate(line)
	// suffix is same length as limit, cutAt would be 0, so we just use maxBytes
	if len(got) > 3 {
		t.Fatalf("expected len <= 3, got %d: %q", len(got), got)
	}
}
