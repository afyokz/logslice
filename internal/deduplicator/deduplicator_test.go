package deduplicator

import (
	"fmt"
	"testing"
)

func TestNew_DisabledWhenWindowZero(t *testing.T) {
	d := New(0)
	if d.Enabled() {
		t.Fatal("expected deduplicator to be disabled for window=0")
	}
}

func TestNew_DisabledWhenWindowNegative(t *testing.T) {
	d := New(-5)
	if d.Enabled() {
		t.Fatal("expected deduplicator to be disabled for negative window")
	}
}

func TestNew_EnabledWhenWindowPositive(t *testing.T) {
	d := New(10)
	if !d.Enabled() {
		t.Fatal("expected deduplicator to be enabled for window=10")
	}
}

func TestIsDuplicate_Disabled_AlwaysFalse(t *testing.T) {
	d := New(0)
	for i := 0; i < 5; i++ {
		if d.IsDuplicate("same line") {
			t.Fatal("disabled deduplicator should never report duplicate")
		}
	}
}

func TestIsDuplicate_FirstOccurrenceNotDuplicate(t *testing.T) {
	d := New(100)
	if d.IsDuplicate("unique line") {
		t.Fatal("first occurrence should not be a duplicate")
	}
}

func TestIsDuplicate_SecondOccurrenceIsDuplicate(t *testing.T) {
	d := New(100)
	d.IsDuplicate("repeated line")
	if !d.IsDuplicate("repeated line") {
		t.Fatal("second occurrence should be a duplicate")
	}
}

func TestIsDuplicate_WindowEviction(t *testing.T) {
	window := 3
	d := New(window)

	// Fill the window with distinct lines
	for i := 0; i < window; i++ {
		d.IsDuplicate(fmt.Sprintf("line-%d", i))
	}

	// Adding one more should evict "line-0"
	d.IsDuplicate("line-new")

	// "line-0" should no longer be considered a duplicate
	if d.IsDuplicate("line-0") {
		t.Fatal("line-0 should have been evicted from the window")
	}
}

func TestReset_ClearsState(t *testing.T) {
	d := New(10)
	d.IsDuplicate("some line")
	d.Reset()
	if d.IsDuplicate("some line") {
		t.Fatal("after Reset, line should not be considered duplicate")
	}
}

func TestIsDuplicate_UniqueLines_NeverDuplicate(t *testing.T) {
	d := New(50)
	for i := 0; i < 50; i++ {
		line := fmt.Sprintf("unique-%d", i)
		if d.IsDuplicate(line) {
			t.Fatalf("line %q should not be a duplicate", line)
		}
	}
}
