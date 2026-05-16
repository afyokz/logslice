package linepadder_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/linepadder"
)

func TestNew_DisabledWhenZero(t *testing.T) {
	p := linepadder.New(0, linepadder.AlignLeft, ' ')
	if p.Enabled() {
		t.Fatal("expected disabled for width=0")
	}
}

func TestNew_DisabledWhenNegative(t *testing.T) {
	p := linepadder.New(-5, linepadder.AlignLeft, ' ')
	if p.Enabled() {
		t.Fatal("expected disabled for negative width")
	}
}

func TestNew_EnabledWhenPositive(t *testing.T) {
	p := linepadder.New(20, linepadder.AlignLeft, ' ')
	if !p.Enabled() {
		t.Fatal("expected enabled for positive width")
	}
	if p.Width() != 20 {
		t.Fatalf("want width 20, got %d", p.Width())
	}
}

func TestPad_Disabled_ReturnsOriginal(t *testing.T) {
	p := linepadder.New(0, linepadder.AlignLeft, ' ')
	got := p.Pad("hello")
	if got != "hello" {
		t.Fatalf("want %q, got %q", "hello", got)
	}
}

func TestPad_AlignLeft_PadsRight(t *testing.T) {
	p := linepadder.New(8, linepadder.AlignLeft, ' ')
	got := p.Pad("hi")
	want := "hi      "
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestPad_AlignRight_PadsLeft(t *testing.T) {
	p := linepadder.New(8, linepadder.AlignRight, ' ')
	got := p.Pad("hi")
	want := "      hi"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestPad_AlignCenter_PadsEvenlyBiasLeft(t *testing.T) {
	p := linepadder.New(8, linepadder.AlignCenter, ' ')
	got := p.Pad("hi")
	want := "   hi   "
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestPad_TruncatesLongLine(t *testing.T) {
	p := linepadder.New(5, linepadder.AlignLeft, ' ')
	got := p.Pad("hello world")
	if got != "hello" {
		t.Fatalf("want %q, got %q", "hello", got)
	}
}

func TestPad_ExactWidth_Unchanged(t *testing.T) {
	p := linepadder.New(5, linepadder.AlignLeft, ' ')
	got := p.Pad("hello")
	if got != "hello" {
		t.Fatalf("want %q, got %q", "hello", got)
	}
}

func TestPad_CustomPadRune(t *testing.T) {
	p := linepadder.New(6, linepadder.AlignLeft, '-')
	got := p.Pad("ab")
	want := "ab----"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestPad_ZeroPadRune_DefaultsToSpace(t *testing.T) {
	p := linepadder.New(5, linepadder.AlignLeft, 0)
	got := p.Pad("x")
	want := "x    "
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}
