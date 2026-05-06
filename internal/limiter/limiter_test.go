package limiter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/limiter"
)

func TestNew_DefaultsToUnlimited(t *testing.T) {
	l := limiter.New(0)
	for i := 0; i < 1000; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow() to return true for unlimited limiter at iteration %d", i)
		}
	}
}

func TestAllow_RespectsMaxLines(t *testing.T) {
	l := limiter.New(3)

	if !l.Allow() {
		t.Fatal("expected true on call 1")
	}
	if !l.Allow() {
		t.Fatal("expected true on call 2")
	}
	if !l.Allow() {
		t.Fatal("expected true on call 3")
	}
	if l.Allow() {
		t.Fatal("expected false on call 4 (limit reached)")
	}
}

func TestCount_TracksAllowedLines(t *testing.T) {
	l := limiter.New(5)
	for i := 0; i < 5; i++ {
		l.Allow()
	}
	if l.Count() != 5 {
		t.Fatalf("expected count 5, got %d", l.Count())
	}
	// Extra call beyond limit should not increment count
	l.Allow()
	if l.Count() != 5 {
		t.Fatalf("expected count to remain 5 after limit, got %d", l.Count())
	}
}

func TestReached_BeforeAndAfterLimit(t *testing.T) {
	l := limiter.New(2)
	if l.Reached() {
		t.Fatal("should not be reached before any calls")
	}
	l.Allow()
	if l.Reached() {
		t.Fatal("should not be reached after 1 of 2")
	}
	l.Allow()
	if !l.Reached() {
		t.Fatal("should be reached after 2 of 2")
	}
}

func TestReached_UnlimitedNeverReaches(t *testing.T) {
	l := limiter.New(0)
	for i := 0; i < 100; i++ {
		l.Allow()
	}
	if l.Reached() {
		t.Fatal("unlimited limiter should never report Reached")
	}
}

func TestReset_ResetsCounter(t *testing.T) {
	l := limiter.New(2)
	l.Allow()
	l.Allow()
	if !l.Reached() {
		t.Fatal("expected Reached before reset")
	}
	l.Reset()
	if l.Count() != 0 {
		t.Fatalf("expected count 0 after reset, got %d", l.Count())
	}
	if l.Reached() {
		t.Fatal("expected not Reached after reset")
	}
	if !l.Allow() {
		t.Fatal("expected Allow to return true after reset")
	}
}

func TestNew_NegativeMaxLines_Unlimited(t *testing.T) {
	l := limiter.New(-5)
	for i := 0; i < 50; i++ {
		if !l.Allow() {
			t.Fatalf("expected unlimited behaviour for negative MaxLines at iteration %d", i)
		}
	}
}
