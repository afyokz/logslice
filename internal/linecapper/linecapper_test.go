package linecapper_test

import (
	"strings"
	"testing"

	"github.com/yourusername/logslice/internal/linecapper"
)

func TestNew_DisabledWhenZero(t *testing.T) {
	c := linecapper.New(0, "...")
	if c.Enabled() {
		t.Fatal("expected disabled capper for maxBytes=0")
	}
}

func TestNew_DisabledWhenNegative(t *testing.T) {
	c := linecapper.New(-5, "...")
	if c.Enabled() {
		t.Fatal("expected disabled capper for negative maxBytes")
	}
}

func TestNew_EnabledWhenPositive(t *testing.T) {
	c := linecapper.New(80, "...")
	if !c.Enabled() {
		t.Fatal("expected enabled capper for positive maxBytes")
	}
	if c.MaxBytes() != 80 {
		t.Fatalf("MaxBytes: got %d, want 80", c.MaxBytes())
	}
}

func TestCap_Disabled_ReturnsOriginal(t *testing.T) {
	c := linecapper.New(0, "...")
	long := strings.Repeat("x", 200)
	if got := c.Cap(long); got != long {
		t.Fatalf("expected original line returned when disabled")
	}
}

func TestCap_ShortLine_Unchanged(t *testing.T) {
	c := linecapper.New(100, "...")
	line := "short line"
	if got := c.Cap(line); got != line {
		t.Fatalf("got %q, want %q", got, line)
	}
}

func TestCap_ExactLength_Unchanged(t *testing.T) {
	c := linecapper.New(10, "...")
	line := "0123456789" // exactly 10 bytes
	if got := c.Cap(line); got != line {
		t.Fatalf("got %q, want %q", got, line)
	}
}

func TestCap_LongLine_TruncatedWithSuffix(t *testing.T) {
	c := linecapper.New(10, "...")
	line := "0123456789ABCDEF"
	got := c.Cap(line)
	if len(got) > 10 {
		t.Fatalf("result length %d exceeds maxBytes 10", len(got))
	}
	if !strings.HasSuffix(got, "...") {
		t.Fatalf("expected suffix '...', got %q", got)
	}
	if got != "0123456..." {
		t.Fatalf("got %q, want %q", got, "0123456...")
	}
}

func TestCap_EmptySuffix_TruncatesCleanly(t *testing.T) {
	c := linecapper.New(5, "")
	line := "ABCDEFGHIJ"
	got := c.Cap(line)
	if got != "ABCDE" {
		t.Fatalf("got %q, want %q", got, "ABCDE")
	}
}

func TestCap_SuffixLargerThanBudget_ReturnsSuffixTruncated(t *testing.T) {
	c := linecapper.New(2, "...")
	line := "hello world"
	got := c.Cap(line)
	if len(got) > 2 {
		t.Fatalf("result length %d exceeds maxBytes 2", len(got))
	}
	if got != ".."[:2] {
		t.Fatalf("got %q, want %q", got, ".."[:2])
	}
}
