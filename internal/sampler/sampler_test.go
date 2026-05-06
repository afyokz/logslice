package sampler_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/sampler"
)

func TestNew_DefaultStep(t *testing.T) {
	s := sampler.New(0)
	if s.Step() != 1 {
		t.Errorf("expected step 1 for zero input, got %d", s.Step())
	}
}

func TestNew_NegativeStep(t *testing.T) {
	s := sampler.New(-5)
	if s.Step() != 1 {
		t.Errorf("expected step 1 for negative input, got %d", s.Step())
	}
}

func TestKeep_StepOne_AlwaysTrue(t *testing.T) {
	s := sampler.New(1)
	for i := 0; i < 10; i++ {
		if !s.Keep() {
			t.Errorf("expected Keep() to return true at iteration %d", i)
		}
	}
}

func TestKeep_StepThree(t *testing.T) {
	s := sampler.New(3)
	expected := []bool{false, false, true, false, false, true}
	for i, want := range expected {
		got := s.Keep()
		if got != want {
			t.Errorf("iteration %d: expected Keep()=%v, got %v", i, want, got)
		}
	}
}

func TestReset_ResetsCounter(t *testing.T) {
	s := sampler.New(3)
	s.Keep()
	s.Keep()
	s.Reset()
	// After reset, next two calls should be false
	if s.Keep() {
		t.Error("expected false after reset on first call")
	}
	if s.Keep() {
		t.Error("expected false after reset on second call")
	}
	if !s.Keep() {
		t.Error("expected true on third call after reset")
	}
}

func TestApply_StepTwo(t *testing.T) {
	s := sampler.New(2)
	input := []string{"a", "b", "c", "d", "e", "f"}
	got := s.Apply(input)
	want := []string{"b", "d", "f"}
	if len(got) != len(want) {
		t.Fatalf("expected %d lines, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: expected %q, got %q", i, want[i], got[i])
		}
	}
}

func TestApply_EmptyInput(t *testing.T) {
	s := sampler.New(2)
	got := s.Apply([]string{})
	if len(got) != 0 {
		t.Errorf("expected empty output, got %v", got)
	}
}

func TestApply_ResetsBeforeEachCall(t *testing.T) {
	s := sampler.New(2)
	input := []string{"a", "b", "c", "d"}
	first := s.Apply(input)
	second := s.Apply(input)
	if len(first) != len(second) {
		t.Errorf("Apply should produce consistent results: first=%v second=%v", first, second)
	}
	for i := range first {
		if first[i] != second[i] {
			t.Errorf("index %d differs: %q vs %q", i, first[i], second[i])
		}
	}
}
