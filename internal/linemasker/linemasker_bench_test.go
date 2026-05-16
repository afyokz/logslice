package linemasker

import "testing"

var sink string

func BenchmarkMask_SinglePattern(b *testing.B) {
	m, err := New([]string{`password=\S+`}, "***")
	if err != nil {
		b.Fatal(err)
	}
	line := "2024-01-15T10:00:00Z user=alice password=hunter2 action=login"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = m.Mask(line)
	}
}

func BenchmarkMask_ThreePatterns(b *testing.B) {
	m, err := New([]string{`password=\S+`, `token=\S+`, `secret=\S+`}, "***")
	if err != nil {
		b.Fatal(err)
	}
	line := "2024-01-15T10:00:00Z user=alice password=hunter2 token=abc123 secret=xyz action=login"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = m.Mask(line)
	}
}

func BenchmarkMask_Disabled(b *testing.B) {
	m, err := New(nil, "")
	if err != nil {
		b.Fatal(err)
	}
	line := "2024-01-15T10:00:00Z user=alice password=hunter2 action=login"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = m.Mask(line)
	}
}
