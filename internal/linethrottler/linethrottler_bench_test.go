package linethrottler

import (
	"testing"
	"time"
)

func BenchmarkApply_Disabled(b *testing.B) {
	th := New(0)
	th.sleepFn = func(time.Duration) {}
	line := "2024-01-01T00:00:00Z INFO benchmark line content here"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = th.Apply(line)
	}
}

func BenchmarkApply_Enabled_MockSleep(b *testing.B) {
	th := New(1 * time.Millisecond)
	th.sleepFn = func(time.Duration) {} // mock to avoid real sleeping
	line := "2024-01-01T00:00:00Z INFO benchmark line content here"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = th.Apply(line)
	}
}

func BenchmarkWait_Disabled(b *testing.B) {
	th := New(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		th.Wait()
	}
}
