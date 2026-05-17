package linescaler_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/linescaler"
)

const benchLine = "2024-01-15T10:00:00Z INFO 987654 request completed successfully"

func BenchmarkScale_Enabled(b *testing.B) {
	s := linescaler.New(0.001, 2, " ")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.Scale(benchLine)
	}
}

func BenchmarkScale_Disabled(b *testing.B) {
	s := linescaler.New(1, 2, " ")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.Scale(benchLine)
	}
}

func BenchmarkScale_CSVMiddleField(b *testing.B) {
	s := linescaler.New(1000, 3, ",")
	line := "ts,INFO,host,250,ok"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.Scale(line)
	}
}
