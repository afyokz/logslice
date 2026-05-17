package linecensor

import (
	"testing"
)

func TestNew_EmptyPatterns_Disabled(t *testing.T) {
	c, err := New(nil, "***")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Enabled() {
		t.Error("expected censor to be disabled")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := New([]string{"[invalid"}, "***")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_ValidPatterns_Enabled(t *testing.T) {
	c, err := New([]string{`\d+`}, "NUM")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !c.Enabled() {
		t.Error("expected censor to be enabled")
	}
}

func TestNew_DefaultPlaceholder(t *testing.T) {
	c, err := New([]string{`\d+`}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Placeholder() != "***" {
		t.Errorf("expected default placeholder '***', got %q", c.Placeholder())
	}
}

func TestApply_Disabled_ReturnsOriginal(t *testing.T) {
	c, _ := New(nil, "***")
	got := c.Apply("hello 123")
	if got != "hello 123" {
		t.Errorf("expected original line, got %q", got)
	}
}

func TestApply_ReplacesMatch(t *testing.T) {
	c, _ := New([]string{`\d+`}, "NUM")
	got := c.Apply("error on line 42")
	if got != "error on line NUM" {
		t.Errorf("unexpected result: %q", got)
	}
}

func TestApply_MultiplePatterns(t *testing.T) {
	c, _ := New([]string{`\d+`, `foo`}, "X")
	got := c.Apply("foo 99 bar")
	if got != "X X bar" {
		t.Errorf("unexpected result: %q", got)
	}
}

func TestContainsAny_Disabled_AlwaysFalse(t *testing.T) {
	c, _ := New(nil, "***")
	if c.ContainsAny("hello 123") {
		t.Error("expected false for disabled censor")
	}
}

func TestContainsAny_Match(t *testing.T) {
	c, _ := New([]string{`secret`}, "***")
	if !c.ContainsAny("my secret key") {
		t.Error("expected true for matching line")
	}
}

func TestContainsAny_NoMatch(t *testing.T) {
	c, _ := New([]string{`secret`}, "***")
	if c.ContainsAny("nothing here") {
		t.Error("expected false for non-matching line")
	}
}

func TestApplyAll_ReturnsNewSlice(t *testing.T) {
	c, _ := New([]string{`\d+`}, "N")
	input := []string{"line 1", "line 2", "no digits"}
	out := c.ApplyAll(input)
	if len(out) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(out))
	}
	if out[0] != "line N" || out[1] != "line N" || out[2] != "no digits" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestString_Disabled(t *testing.T) {
	c, _ := New(nil, "")
	if c.String() != "Censor(disabled)" {
		t.Errorf("unexpected String(): %s", c.String())
	}
}
