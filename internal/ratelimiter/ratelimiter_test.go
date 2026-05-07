package ratelimiter

import (
	"testing"
	"time"
)

func TestNew_DisabledWhenZero(t *testing.T) {
	rl := New(0)
	if rl.Enabled() {
		t.Fatal("expected limiter to be disabled for rate=0")
	}
}

func TestNew_DisabledWhenNegative(t *testing.T) {
	rl := New(-5)
	if rl.Enabled() {
		t.Fatal("expected limiter to be disabled for negative rate")
	}
}

func TestNew_EnabledWhenPositive(t *testing.T) {
	rl := New(10)
	if !rl.Enabled() {
		t.Fatal("expected limiter to be enabled for rate=10")
	}
	if rl.Rate() != 10 {
		t.Fatalf("expected rate 10, got %d", rl.Rate())
	}
}

func TestAllow_Disabled_AlwaysTrue(t *testing.T) {
	rl := New(0)
	for i := 0; i < 1000; i++ {
		if !rl.Allow() {
			t.Fatalf("disabled limiter denied line %d", i)
		}
	}
}

func TestAllow_RespectsRate(t *testing.T) {
	fixed := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	rl := New(3)
	rl.now = func() time.Time { return fixed }
	rl.lastFill = fixed
	rl.bucket = 3

	// First 3 should be allowed.
	for i := 0; i < 3; i++ {
		if !rl.Allow() {
			t.Fatalf("expected Allow()=true for call %d", i)
		}
	}
	// 4th should be denied (bucket exhausted, no time elapsed).
	if rl.Allow() {
		t.Fatal("expected Allow()=false after bucket exhausted")
	}
}

func TestAllow_RefillsAfterOneSecond(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	rl := New(5)
	rl.now = func() time.Time { return base }
	rl.lastFill = base
	rl.bucket = 0 // Start exhausted.

	// Still exhausted — no time passed.
	if rl.Allow() {
		t.Fatal("expected denial when bucket is empty")
	}

	// Advance clock by 1 second — bucket should refill to 5.
	advanced := base.Add(time.Second)
	rl.now = func() time.Time { return advanced }

	for i := 0; i < 5; i++ {
		if !rl.Allow() {
			t.Fatalf("expected Allow()=true after refill, call %d", i)
		}
	}
	if rl.Allow() {
		t.Fatal("expected denial after new bucket exhausted")
	}
}

func TestAllow_BucketCappedAtRate(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	rl := New(4)
	rl.now = func() time.Time { return base }
	rl.lastFill = base
	rl.bucket = 0

	// Advance by 10 seconds — bucket should be capped at rate (4), not 40.
	advanced := base.Add(10 * time.Second)
	rl.now = func() time.Time { return advanced }

	allowed := 0
	for i := 0; i < 10; i++ {
		if rl.Allow() {
			allowed++
		}
	}
	if allowed != 4 {
		t.Fatalf("expected bucket capped at 4, got %d allowed", allowed)
	}
}
