package linewindow_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/linewindow"
)

func TestNew_DisabledWhenZero(t *testing.T) {
	w := linewindow.New(0)
	if w.Enabled() {
		t.Fatal("expected disabled window")
	}
}

func TestNew_DisabledWhenNegative(t *testing.T) {
	w := linewindow.New(-5)
	if w.Enabled() {
		t.Fatal("expected disabled window for negative size")
	}
}

func TestNew_EnabledWhenPositive(t *testing.T) {
	w := linewindow.New(3)
	if !w.Enabled() {
		t.Fatal("expected enabled window")
	}
}

func TestPush_Disabled_NoOp(t *testing.T) {
	w := linewindow.New(0)
	w.Push("line")
	if w.Len() != 0 {
		t.Fatalf("expected len 0, got %d", w.Len())
	}
}

func TestLines_Empty_ReturnsNil(t *testing.T) {
	w := linewindow.New(3)
	if w.Lines() != nil {
		t.Fatal("expected nil for empty window")
	}
}

func TestPush_UnderCapacity_PreservesOrder(t *testing.T) {
	w := linewindow.New(5)
	w.Push("a")
	w.Push("b")
	w.Push("c")
	got := w.Lines()
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("len mismatch: got %d want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q want %q", i, got[i], want[i])
		}
	}
}

func TestPush_OverCapacity_EvictsOldest(t *testing.T) {
	w := linewindow.New(3)
	for _, l := range []string{"a", "b", "c", "d", "e"} {
		w.Push(l)
	}
	got := w.Lines()
	want := []string{"c", "d", "e"}
	if len(got) != len(want) {
		t.Fatalf("len mismatch: got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q want %q", i, got[i], want[i])
		}
	}
}

func TestReset_ClearsBuffer(t *testing.T) {
	w := linewindow.New(3)
	w.Push("x")
	w.Push("y")
	w.Reset()
	if w.Len() != 0 {
		t.Fatalf("expected len 0 after reset, got %d", w.Len())
	}
	if w.Lines() != nil {
		t.Fatal("expected nil lines after reset")
	}
}

func TestLines_ReturnsCopy(t *testing.T) {
	w := linewindow.New(3)
	w.Push("a")
	lines := w.Lines()
	lines[0] = "mutated"
	got := w.Lines()
	if got[0] != "a" {
		t.Errorf("Lines should return a copy; got %q", got[0])
	}
}
