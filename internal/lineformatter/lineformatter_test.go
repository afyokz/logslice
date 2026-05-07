package lineformatter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/lineformatter"
)

func TestNew_NoNumbering_NoPrefix(t *testing.T) {
	f := lineformatter.New(false, "")
	if f.Numbered() {
		t.Fatal("expected numbered=false")
	}
	if f.Prefix() != "" {
		t.Fatalf("expected empty prefix, got %q", f.Prefix())
	}
}

func TestNew_NumberedEnabled(t *testing.T) {
	f := lineformatter.New(true, "")
	if !f.Numbered() {
		t.Fatal("expected numbered=true")
	}
}

func TestNew_PrefixStored(t *testing.T) {
	f := lineformatter.New(false, ">> ")
	if f.Prefix() != ">> " {
		t.Fatalf("expected prefix '>> ', got %q", f.Prefix())
	}
}

func TestFormat_PlainLine_Unchanged(t *testing.T) {
	f := lineformatter.New(false, "")
	got := f.Format(1, "hello world")
	if got != "hello world" {
		t.Fatalf("expected 'hello world', got %q", got)
	}
}

func TestFormat_Numbered_PrependsN(t *testing.T) {
	f := lineformatter.New(true, "")
	got := f.Format(3, "some log line")
	want := "3 some log line"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestFormat_Prefix_PrependsString(t *testing.T) {
	f := lineformatter.New(false, "[INFO] ")
	got := f.Format(1, "message")
	want := "[INFO] message"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestFormat_NumberedAndPrefix(t *testing.T) {
	f := lineformatter.New(true, ">> ")
	got := f.Format(7, "entry")
	want := "7 >> entry"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestFormat_NumberStartsAtOne(t *testing.T) {
	f := lineformatter.New(true, "")
	first := f.Format(1, "first")
	if first != "1 first" {
		t.Fatalf("expected '1 first', got %q", first)
	}
	second := f.Format(2, "second")
	if second != "2 second" {
		t.Fatalf("expected '2 second', got %q", second)
	}
}

func TestFormat_EmptyLine_Numbered(t *testing.T) {
	f := lineformatter.New(true, "")
	got := f.Format(5, "")
	if got != "5 " {
		t.Fatalf("expected '5 ', got %q", got)
	}
}
