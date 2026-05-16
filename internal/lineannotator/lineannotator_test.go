package lineannotator

import (
	"testing"
)

func TestNew_EmptyTag_Disabled(t *testing.T) {
	a := New("", "")
	if a.Enabled() {
		t.Fatal("expected disabled annotator for empty tag")
	}
}

func TestNew_WithTag_Enabled(t *testing.T) {
	a := New("SRC", "")
	if !a.Enabled() {
		t.Fatal("expected enabled annotator for non-empty tag")
	}
}

func TestNew_DefaultSeparator(t *testing.T) {
	a := New("APP", "")
	if a.sep != ": " {
		t.Fatalf("expected default sep \": \", got %q", a.sep)
	}
}

func TestNew_CustomSeparator(t *testing.T) {
	a := New("APP", " | ")
	if a.sep != " | " {
		t.Fatalf("expected sep \" | \", got %q", a.sep)
	}
}

func TestTag_ReturnsConfiguredTag(t *testing.T) {
	a := New("myservice", "")
	if a.Tag() != "myservice" {
		t.Fatalf("expected tag %q, got %q", "myservice", a.Tag())
	}
}

func TestAnnotate_Disabled_ReturnsOriginal(t *testing.T) {
	a := New("", "")
	line := "some log line"
	if got := a.Annotate(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestAnnotate_PrependsTagnSep(t *testing.T) {
	a := New("web", ": ")
	got := a.Annotate("request received")
	want := "web: request received"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestAnnotate_CustomSeparator(t *testing.T) {
	a := New("db", " >> ")
	got := a.Annotate("query executed")
	want := "db >> query executed"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestAnnotate_EmptyLine(t *testing.T) {
	a := New("svc", ": ")
	got := a.Annotate("")
	want := "svc: "
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
