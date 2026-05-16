package lineskipper_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/lineskipper"
)

func TestNew_DisabledWhenZero(t *testing.T) {
	s := lineskipper.New(0)
	if s.Enabled() {
		t.Fatal("expected disabled when skip=0")
	}
}

func TestNew_DisabledWhenNegative(t *testing.T) {
	s := lineskipper.New(-5)
	if s.Enabled() {
		t.Fatal("expected disabled when skip<0")
	}
}

func TestNew_EnabledWhenPositive(t *testing.T) {
	s := lineskipper.New(3)
	if !s.Enabled() {
		t.Fatal("expected enabled when skip>0")
	}
	if s.Skip() != 3 {
		t.Fatalf("expected skip=3, got %d", s.Skip())
	}
}

func TestKeep_Disabled_AlwaysTrue(t *testing.T) {
	s := lineskipper.New(0)
	for i := 0; i < 5; i++ {
		if !s.Keep("line") {
			t.Fatalf("expected Keep=true on disabled skipper, iteration %d", i)
		}
	}
}

func TestKeep_DropsFirstNLines(t *testing.T) {
	s := lineskipper.New(3)
	results := make([]bool, 6)
	for i := range results {
		results[i] = s.Keep("line")
	}
	expected := []bool{false, false, false, true, true, true}
	for i, got := range results {
		if got != expected[i] {
			t.Errorf("line %d: expected Keep=%v, got %v", i, expected[i], got)
		}
	}
}

func TestCount_TracksDroppedLines(t *testing.T) {
	s := lineskipper.New(2)
	if s.Count() != 0 {
		t.Fatal("expected count=0 initially")
	}
	s.Keep("a")
	s.Keep("b")
	s.Keep("c") // this one is kept, count stays at 2
	if s.Count() != 2 {
		t.Fatalf("expected count=2, got %d", s.Count())
	}
}

func TestReset_ResetsCounter(t *testing.T) {
	s := lineskipper.New(2)
	s.Keep("a")
	s.Keep("b")
	if s.Count() != 2 {
		t.Fatal("pre-reset: expected count=2")
	}
	s.Reset()
	if s.Count() != 0 {
		t.Fatalf("post-reset: expected count=0, got %d", s.Count())
	}
	// After reset the first two lines should be dropped again.
	if s.Keep("x") {
		t.Fatal("expected first line after reset to be dropped")
	}
}

func TestKeep_SkipOne(t *testing.T) {
	s := lineskipper.New(1)
	if s.Keep("header") {
		t.Fatal("expected first line to be dropped")
	}
	if !s.Keep("data") {
		t.Fatal("expected second line to be kept")
	}
}
