package linechunker

import "testing"

func BenchmarkFeed_ChunkSize10(b *testing.B) {
	c := New(10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if ch := c.Feed("2024-01-01T00:00:00Z some log line content here"); ch != nil {
			_ = ch
		}
	}
	_ = c.Flush()
}

func BenchmarkFeed_ChunkSize100(b *testing.B) {
	c := New(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if ch := c.Feed("2024-01-01T00:00:00Z some log line content here"); ch != nil {
			_ = ch
		}
	}
	_ = c.Flush()
}

func BenchmarkFeed_Disabled(b *testing.B) {
	c := New(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.Feed("2024-01-01T00:00:00Z some log line content here")
	}
}
