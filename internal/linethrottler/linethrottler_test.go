package linethrottler

import (
	"testing"
	"time"
)

func TestNew_DisabledWhenZero(t *testing.T) {
	th := New(0)
	if th.Enabled() {
		t.Fatal("expected disabled for zero delay")
	}
}

func TestNew_DisabledWhenNegative(t *testing.T) {
	th := New(-1 * time.Millisecond)
	if th.Enabled() {
		t.Fatal("expected disabled for negative delay")
	}
}

func TestNew_EnabledWhenPositive(t *testing.T) {
	th := New(10 * time.Millisecond)
	if !th.Enabled() {
		t.Fatal("expected enabled for positive delay")
	}
}

func TestDelay_ReturnsConfiguredValue(t *testing.T) {
	d := 42 * time.Millisecond
	th := New(d)
	if th.Delay() != d {
		t.Fatalf("want %v, got %v", d, th.Delay())
	}
}

func TestWait_Disabled_DoesNotSleep(t *testing.T) {
	th := New(0)
	called := false
	th.sleepFn = func(d time.Duration) { called = true }
	th.Wait()
	if called {
		t.Fatal("sleepFn should not be called when disabled")
	}
}

func TestWait_Enabled_Sleeps(t *testing.T) {
	th := New(5 * time.Millisecond)
	var slept time.Duration
	th.sleepFn = func(d time.Duration) { slept = d }
	th.Wait()
	if slept != 5*time.Millisecond {
		t.Fatalf("want 5ms sleep, got %v", slept)
	}
}

func TestApply_Disabled_ReturnsOriginal(t *testing.T) {
	th := New(0)
	th.sleepFn = func(d time.Duration) {}
	line := "hello world"
	if got := th.Apply(line); got != line {
		t.Fatalf("want %q, got %q", line, got)
	}
}

func TestApply_Enabled_ReturnsOriginalAndSleeps(t *testing.T) {
	th := New(1 * time.Millisecond)
	sleepCalls := 0
	th.sleepFn = func(d time.Duration) { sleepCalls++ }
	lines := []string{"line1", "line2", "line3"}
	for _, l := range lines {
		got := th.Apply(l)
		if got != l {
			t.Fatalf("want %q, got %q", l, got)
		}
	}
	if sleepCalls != len(lines) {
		t.Fatalf("want %d sleep calls, got %d", len(lines), sleepCalls)
	}
}
