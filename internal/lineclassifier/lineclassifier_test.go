package lineclassifier_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/lineclassifier"
)

func TestNew_EmptyPatterns_Disabled(t *testing.T) {
	c, err := lineclassifier.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Enabled() {
		t.Fatal("expected classifier to be disabled")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := lineclassifier.New(map[string]string{"bad": "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_ValidPatterns_Enabled(t *testing.T) {
	c, err := lineclassifier.New(map[string]string{"error": `(?i)error`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !c.Enabled() {
		t.Fatal("expected classifier to be enabled")
	}
}

func TestClassify_Disabled_ReturnsEmpty(t *testing.T) {
	c, _ := lineclassifier.New(nil)
	if got := c.Classify("ERROR: something broke"); got != "" {
		t.Fatalf("expected empty label, got %q", got)
	}
}

func TestClassify_MatchesCorrectLabel(t *testing.T) {
	c, err := lineclassifier.New(map[string]string{
		"error": `(?i)\berror\b`,
		"warn":  `(?i)\bwarn\b`,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tests := []struct {
		line  string
		want  string
	}{
		{"2024-01-01 ERROR: disk full", "error"},
		{"2024-01-01 WARN: low memory", "warn"},
		{"2024-01-01 INFO: all good", ""},
	}
	for _, tt := range tests {
		got := c.Classify(tt.line)
		if got != tt.want {
			t.Errorf("Classify(%q) = %q, want %q", tt.line, got, tt.want)
		}
	}
}

func TestClassify_NoMatch_ReturnsEmpty(t *testing.T) {
	c, _ := lineclassifier.New(map[string]string{"error": `(?i)error`})
	if got := c.Classify("INFO: everything is fine"); got != "" {
		t.Fatalf("expected empty label, got %q", got)
	}
}

func TestLabels_ReturnsAllLabels(t *testing.T) {
	c, err := lineclassifier.New(map[string]string{
		"error": `error`,
		"debug": `debug`,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	labels := c.Labels()
	if len(labels) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(labels))
	}
}

func TestLabels_Disabled_ReturnsEmpty(t *testing.T) {
	c, _ := lineclassifier.New(nil)
	if got := c.Labels(); len(got) != 0 {
		t.Fatalf("expected no labels, got %v", got)
	}
}
