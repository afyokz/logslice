package linefilter_test

import (
	"strings"
	"testing"

	"github.com/example/logslice/internal/linefilter"
)

func TestNew_DisabledWhenBothZero(t *testing.T) {
	f := linefilter.New(0, 0)
	if f.Enabled() {
		t.Fatal("expected filter to be disabled when both bounds are zero")
	}
}

func TestNew_EnabledWhenMinSet(t *testing.T) {
	f := linefilter.New(5, 0)
	if !f.Enabled() {
		t.Fatal("expected filter to be enabled when minLen > 0")
	}
}

func TestNew_EnabledWhenMaxSet(t *testing.T) {
	f := linefilter.New(0, 100)
	if !f.Enabled() {
		t.Fatal("expected filter to be enabled when maxLen > 0")
	}
}

func TestKeep_Disabled_AlwaysTrue(t *testing.T) {
	f := linefilter.New(0, 0)
	for _, line := range []string{"", "x", strings.Repeat("a", 1000)} {
		if !f.Keep(line) {
			t.Errorf("disabled filter: Keep(%q) = false, want true", line)
		}
	}
}

func TestKeep_MinLen_RejectsShorterLines(t *testing.T) {
	f := linefilter.New(10, 0)
	if f.Keep("short") {
		t.Error("expected Keep to return false for line shorter than minLen")
	}
	if !f.Keep(strings.Repeat("a", 10)) {
		t.Error("expected Keep to return true for line equal to minLen")
	}
	if !f.Keep(strings.Repeat("a", 20)) {
		t.Error("expected Keep to return true for line longer than minLen")
	}
}

func TestKeep_MaxLen_RejectsLongerLines(t *testing.T) {
	f := linefilter.New(0, 20)
	if f.Keep(strings.Repeat("a", 21)) {
		t.Error("expected Keep to return false for line longer than maxLen")
	}
	if !f.Keep(strings.Repeat("a", 20)) {
		t.Error("expected Keep to return true for line equal to maxLen")
	}
	if !f.Keep("hi") {
		t.Error("expected Keep to return true for line shorter than maxLen")
	}
}

func TestKeep_BothBounds(t *testing.T) {
	f := linefilter.New(5, 15)
	cases := []struct {
		line string
		want bool
	}{
		{"abc", false},
		{"abcde", true},
		{strings.Repeat("x", 15), true},
		{strings.Repeat("x", 16), false},
	}
	for _, tc := range cases {
		got := f.Keep(tc.line)
		if got != tc.want {
			t.Errorf("Keep(%q) = %v, want %v", tc.line, got, tc.want)
		}
	}
}

func TestStats_ReturnsConfiguredBounds(t *testing.T) {
	f := linefilter.New(3, 42)
	min, max := f.Stats()
	if min != 3 || max != 42 {
		t.Errorf("Stats() = (%d, %d), want (3, 42)", min, max)
	}
}
