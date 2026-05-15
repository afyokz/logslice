package linecolumnextractor_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/linecolumnextractor"
)

var sink map[string]string

func BenchmarkExtract_ThreeColumns(b *testing.B) {
	e := linecolumnextractor.New(" ", []string{"timestamp", "level", "message"})
	line := "2024-06-01T12:00:00Z ERROR something went wrong in the subsystem"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = e.Extract(line)
	}
}

func BenchmarkExtract_TenColumns(b *testing.B) {
	cols := []string{"c1", "c2", "c3", "c4", "c5", "c6", "c7", "c8", "c9", "c10"}
	e := linecolumnextractor.New("|", cols)
	line := "a|b|c|d|e|f|g|h|i|j"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = e.Extract(line)
	}
}

func BenchmarkExtract_Disabled(b *testing.B) {
	e := linecolumnextractor.New(" ", nil)
	line := "2024-06-01T12:00:00Z ERROR something went wrong"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = e.Extract(line)
	}
}
