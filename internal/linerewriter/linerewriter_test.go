package linerewriter

import (
	"testing"
)

func TestNew_EmptyRules_Disabled(t *testing.T) {
	rw, err := New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rw.Enabled() {
		t.Fatal("expected rewriter to be disabled")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{Pattern: "[invalid", Replacement: "x"}})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_ValidRules_Enabled(t *testing.T) {
	rw, err := New([]Rule{{Pattern: `\d+`, Replacement: "NUM"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rw.Enabled() {
		t.Fatal("expected rewriter to be enabled")
	}
}

func TestRewrite_Disabled_ReturnsOriginal(t *testing.T) {
	rw, _ := New(nil)
	got := rw.Rewrite("hello world")
	if got != "hello world" {
		t.Fatalf("expected original line, got %q", got)
	}
}

func TestRewrite_SingleRule_ReplacesMatch(t *testing.T) {
	rw, _ := New([]Rule{{Pattern: `\d+`, Replacement: "NUM"}})
	got := rw.Rewrite("error at line 42")
	want := "error at line NUM"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestRewrite_MultipleRules_AppliedInOrder(t *testing.T) {
	rw, _ := New([]Rule{
		{Pattern: `foo`, Replacement: "bar"},
		{Pattern: `bar`, Replacement: "baz"},
	})
	// first rule turns foo->bar, second turns bar->baz
	got := rw.Rewrite("foo")
	if got != "baz" {
		t.Fatalf("want %q, got %q", "baz", got)
	}
}

func TestRewrite_NoMatch_ReturnsOriginal(t *testing.T) {
	rw, _ := New([]Rule{{Pattern: `xyz`, Replacement: "ABC"}})
	got := rw.Rewrite("hello world")
	if got != "hello world" {
		t.Fatalf("expected original line, got %q", got)
	}
}

func TestRewrite_EmptyLine_ReturnsEmpty(t *testing.T) {
	rw, _ := New([]Rule{{Pattern: `\d+`, Replacement: "NUM"}})
	got := rw.Rewrite("")
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestRewrite_GlobalReplacement(t *testing.T) {
	rw, _ := New([]Rule{{Pattern: `\d`, Replacement: "#"}})
	got := rw.Rewrite("a1b2c3")
	want := "a#b#c#"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}
