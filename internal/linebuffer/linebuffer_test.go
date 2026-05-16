package linebuffer_test

import (
	"testing"

	"github.com/logslice/logslice/internal/linebuffer"
)

func TestNew_DisabledWhenZero(t *testing.T) {
	b := linebuffer.New(0)
	if b.Enabled() {
		t.Fatal("expected buffer to be disabled")
	}
}

func TestNew_DisabledWhenNegative(t *testing.T) {
	b := linebuffer.New(-5)
	if b.Enabled() {
		t.Fatal("expected buffer to be disabled")
	}
}

func TestNew_EnabledWhenPositive(t *testing.T) {
	b := linebuffer.New(3)
	if !b.Enabled() {
		t.Fatal("expected buffer to be enabled")
	}
}

func TestPush_Disabled_NoOp(t *testing.T) {
	b := linebuffer.New(0)
	b.Push("line")
	if b.Len() != 0 {
		t.Fatalf("expected len 0, got %d", b.Len())
	}
}

func TestLines_Empty_ReturnsNil(t *testing.T) {
	b := linebuffer.New(4)
	if b.Lines() != nil {
		t.Fatal("expected nil for empty buffer")
	}
}

func TestPush_LessThanCap_OrderPreserved(t *testing.T) {
	b := linebuffer.New(5)
	b.Push("a")
	b.Push("b")
	b.Push("c")
	lines := b.Lines()
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	expected := []string{"a", "b", "c"}
	for i, want := range expected {
		if lines[i] != want {
			t.Errorf("lines[%d] = %q, want %q", i, lines[i], want)
		}
	}
}

func TestPush_OverCap_EvictsOldest(t *testing.T) {
	b := linebuffer.New(3)
	for _, l := range []string{"a", "b", "c", "d", "e"} {
		b.Push(l)
	}
	lines := b.Lines()
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	expected := []string{"c", "d", "e"}
	for i, want := range expected {
		if lines[i] != want {
			t.Errorf("lines[%d] = %q, want %q", i, lines[i], want)
		}
	}
}

func TestReset_ClearsBuffer(t *testing.T) {
	b := linebuffer.New(4)
	b.Push("x")
	b.Push("y")
	b.Reset()
	if b.Len() != 0 {
		t.Fatalf("expected len 0 after reset, got %d", b.Len())
	}
	if b.Lines() != nil {
		t.Fatal("expected nil lines after reset")
	}
}

func TestLen_TracksSize(t *testing.T) {
	b := linebuffer.New(10)
	for i := 0; i < 7; i++ {
		b.Push("line")
	}
	if b.Len() != 7 {
		t.Fatalf("expected len 7, got %d", b.Len())
	}
}
