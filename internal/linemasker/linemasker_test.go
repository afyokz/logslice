package linemasker

import (
	"testing"
)

func TestNew_EmptyPatterns_Disabled(t *testing.T) {
	m, err := New(nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Enabled() {
		t.Error("expected masker to be disabled when no patterns provided")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := New([]string{"[invalid"}, "***")
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestNew_ValidPatterns_Enabled(t *testing.T) {
	m, err := New([]string{`\d+`}, "NUM")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Enabled() {
		t.Error("expected masker to be enabled")
	}
}

func TestNew_DefaultReplacement(t *testing.T) {
	m, err := New([]string{`\d+`}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Replacement() != "***" {
		t.Errorf("expected default replacement '***', got %q", m.Replacement())
	}
}

func TestMask_Disabled_ReturnsOriginal(t *testing.T) {
	m, _ := New(nil, "")
	line := "user=admin password=secret123"
	if got := m.Mask(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestMask_SinglePattern(t *testing.T) {
	m, err := New([]string{`password=\S+`}, "password=[REDACTED]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := "user=admin password=secret123 role=reader"
	want := "user=admin password=[REDACTED] role=reader"
	if got := m.Mask(input); got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestMask_MultiplePatterns(t *testing.T) {
	m, err := New([]string{`password=\S+`, `token=\S+`}, "***")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := "password=abc token=xyz user=bob"
	want := "*** *** user=bob"
	if got := m.Mask(input); got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestMask_NoMatch_ReturnsOriginal(t *testing.T) {
	m, err := New([]string{`secret=\S+`}, "***")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := "user=admin role=reader"
	if got := m.Mask(line); got != line {
		t.Errorf("expected unchanged line %q, got %q", line, got)
	}
}

func TestMask_AllOccurrencesReplaced(t *testing.T) {
	m, err := New([]string{`\d+`}, "N")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := "error 404 at line 23 code 500"
	want := "error N at line N code N"
	if got := m.Mask(input); got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
