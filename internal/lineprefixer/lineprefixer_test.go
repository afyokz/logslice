package lineprefixer_test

import (
	"testing"

	"github.com/logslice/logslice/internal/lineprefixer"
)

func TestNew_EmptyPrefix_Disabled(t *testing.T) {
	p := lineprefixer.New("", "")
	if p.Enabled() {
		t.Fatal("expected disabled when prefix is empty")
	}
}

func TestNew_WithPrefix_Enabled(t *testing.T) {
	p := lineprefixer.New("[INFO]", "")
	if !p.Enabled() {
		t.Fatal("expected enabled when prefix is non-empty")
	}
}

func TestNew_EmptySeparator_DefaultsToSpace(t *testing.T) {
	p := lineprefixer.New("[INFO]", "")
	if p.Separator() != " " {
		t.Fatalf("expected separator ' ', got %q", p.Separator())
	}
}

func TestNew_CustomSeparator_Stored(t *testing.T) {
	p := lineprefixer.New("[INFO]", ": ")
	if p.Separator() != ": " {
		t.Fatalf("expected separator ': ', got %q", p.Separator())
	}
}

func TestNew_PrefixStored(t *testing.T) {
	p := lineprefixer.New("[WARN]", "")
	if p.Prefix() != "[WARN]" {
		t.Fatalf("expected prefix '[WARN]', got %q", p.Prefix())
	}
}

func TestFormat_Disabled_ReturnsOriginal(t *testing.T) {
	p := lineprefixer.New("", "")
	got := p.Format("some log line")
	if got != "some log line" {
		t.Fatalf("expected original line, got %q", got)
	}
}

func TestFormat_PrependsPrefixAndSeparator(t *testing.T) {
	p := lineprefixer.New("[INFO]", " ")
	got := p.Format("server started")
	want := "[INFO] server started"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestFormat_CustomSeparator(t *testing.T) {
	p := lineprefixer.New(">>>", ": ")
	got := p.Format("hello world")
	want := ">>>: hello world"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestFormat_EmptyLine_OnlyPrefixAndSeparator(t *testing.T) {
	p := lineprefixer.New("[TAG]", "-")
	got := p.Format("")
	want := "[TAG]-"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestFormat_MultipleCallsAreIdempotent(t *testing.T) {
	p := lineprefixer.New("X", ":")
	line := "test"
	for i := 0; i < 3; i++ {
		got := p.Format(line)
		if got != "X:test" {
			t.Fatalf("call %d: expected %q, got %q", i+1, "X:test", got)
		}
	}
}
