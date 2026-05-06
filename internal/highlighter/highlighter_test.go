package highlighter

import (
	"strings"
	"testing"
)

func TestNew_DisabledHighlighter(t *testing.T) {
	h, err := New(false, "", []string{"error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.Enabled() {
		t.Error("expected highlighter to be disabled")
	}
}

func TestNew_EnabledNoTerms(t *testing.T) {
	h, err := New(true, "yellow", []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.PatternCount() != 0 {
		t.Errorf("expected 0 patterns, got %d", h.PatternCount())
	}
}

func TestNew_EnabledWithTerms(t *testing.T) {
	h, err := New(true, "red", []string{"error", "warn"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.PatternCount() != 2 {
		t.Errorf("expected 2 patterns, got %d", h.PatternCount())
	}
}

func TestHighlight_Disabled_ReturnsOriginal(t *testing.T) {
	h, _ := New(false, "", []string{"error"})
	line := "this is an error line"
	if got := h.Highlight(line); got != line {
		t.Errorf("expected original line, got %q", got)
	}
}

func TestHighlight_NoPatterns_ReturnsOriginal(t *testing.T) {
	h, _ := New(true, "", []string{})
	line := "nothing to highlight"
	if got := h.Highlight(line); got != line {
		t.Errorf("expected original line, got %q", got)
	}
}

func TestHighlight_MatchesAndWraps(t *testing.T) {
	h, err := New(true, "yellow", []string{"error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := "2024-01-01 ERROR something failed"
	result := h.Highlight(line)
	if !strings.Contains(result, ansiYellow) {
		t.Error("expected ANSI yellow code in output")
	}
	if !strings.Contains(result, ansiReset) {
		t.Error("expected ANSI reset code in output")
	}
	if !strings.Contains(result, "ERROR") {
		t.Error("expected original text preserved in output")
	}
}

func TestHighlight_CaseInsensitive(t *testing.T) {
	h, _ := New(true, "red", []string{"warn"})
	result := h.Highlight("2024-01-01 WARN disk usage high")
	if !strings.Contains(result, ansiRed) {
		t.Error("expected red ANSI code for case-insensitive match")
	}
}

func TestHighlight_MultipleTerms(t *testing.T) {
	h, _ := New(true, "", []string{"error", "failed"})
	line := "error: connection failed"
	result := h.Highlight(line)
	count := strings.Count(result, ansiReset)
	if count < 2 {
		t.Errorf("expected at least 2 highlighted segments, resets found: %d", count)
	}
}

func TestResolveColor_Defaults(t *testing.T) {
	if resolveColor("") != ansiYellow {
		t.Error("empty color should default to yellow")
	}
	if resolveColor("unknown") != ansiYellow {
		t.Error("unknown color should default to yellow")
	}
	if resolveColor("red") != ansiRed {
		t.Error("red color should resolve to ansiRed")
	}
}
