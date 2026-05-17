package linenormalizer

import (
	"testing"
)

func TestNew_DisabledWhenFalse(t *testing.T) {
	n := New(false, false)
	if n.Enabled() {
		t.Fatal("expected disabled")
	}
}

func TestNew_EnabledWhenTrue(t *testing.T) {
	n := New(true, false)
	if !n.Enabled() {
		t.Fatal("expected enabled")
	}
}

func TestNew_TrimStored(t *testing.T) {
	n := New(true, true)
	if !n.Trim() {
		t.Fatal("expected trim=true")
	}
}

func TestNormalize_Disabled_ReturnsOriginal(t *testing.T) {
	n := New(false, false)
	input := "hello   world\t!"
	if got := n.Normalize(input); got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

func TestNormalize_CollapsesSpaces(t *testing.T) {
	n := New(true, false)
	got := n.Normalize("hello   world")
	want := "hello world"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestNormalize_CollapsesTabs(t *testing.T) {
	n := New(true, false)
	got := n.Normalize("col1\t\tcol2")
	want := "col1 col2"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestNormalize_MixedWhitespace(t *testing.T) {
	n := New(true, false)
	got := n.Normalize("a \t b")
	want := "a b"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestNormalize_TrimEnabled_RemovesEdges(t *testing.T) {
	n := New(true, true)
	got := n.Normalize("  hello   world  ")
	want := "hello world"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestNormalize_TrimDisabled_PreservesEdges(t *testing.T) {
	n := New(true, false)
	got := n.Normalize("  hello   world  ")
	want := " hello world "
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestNormalize_EmptyString(t *testing.T) {
	n := New(true, true)
	if got := n.Normalize(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}
