package scanner_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/scanner"
)

func TestScan_FiltersByTimeRange(t *testing.T) {
	const layout = "2006-01-02T15:04:05"

	input := strings.Join([]string{
		"2024-03-01T10:00:00 INFO  startup complete\n",
		"2024-03-01T10:05:00 DEBUG request received\n",
		"2024-03-01T10:10:00 ERROR disk full\n",
		"2024-03-01T10:15:00 INFO  shutdown\n",
	}, "")

	start, _ := time.Parse(layout, "2024-03-01T10:04:00")
	end, _ := time.Parse(layout, "2024-03-01T10:11:00")

	opts := scanner.Options{
		Format: layout,
		Start:  start,
		End:    end,
	}

	s := scanner.New(strings.NewReader(input), opts)
	var out bytes.Buffer
	n, err := s.Scan(&out)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 lines written, got %d", n)
	}
	if !strings.Contains(out.String(), "DEBUG request received") {
		t.Error("expected DEBUG line in output")
	}
	if !strings.Contains(out.String(), "ERROR disk full") {
		t.Error("expected ERROR line in output")
	}
	if strings.Contains(out.String(), "startup") {
		t.Error("startup line should be excluded")
	}
	if strings.Contains(out.String(), "shutdown") {
		t.Error("shutdown line should be excluded")
	}
}

func TestScan_EmptyInput(t *testing.T) {
	opts := scanner.Options{
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
	}
	s := scanner.New(strings.NewReader(""), opts)
	var out bytes.Buffer
	n, err := s.Scan(&out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 lines, got %d", n)
	}
}

func TestScan_UnparsableLines_Skipped(t *testing.T) {
	input := "not a timestamp line\nanother bad line\n"
	opts := scanner.Options{
		Start: time.Time{},
		End:   time.Now().Add(24 * time.Hour),
	}
	s := scanner.New(strings.NewReader(input), opts)
	var out bytes.Buffer
	n, err := s.Scan(&out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 lines written for unparsable input, got %d", n)
	}
}
