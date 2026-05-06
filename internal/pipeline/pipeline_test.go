package pipeline_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/logslice/logslice/internal/pipeline"
)

const sampleLog = `2024-01-15T10:00:00Z INFO  service started
2024-01-15T10:01:00Z DEBUG health check ok
2024-01-15T10:02:00Z ERROR disk full
2024-01-15T10:03:00Z INFO  backup completed
2024-01-15T10:04:00Z WARN  memory high
`

func TestRun_BasicSlice(t *testing.T) {
	var out bytes.Buffer
	res, err := pipeline.Run(pipeline.Config{
		Reader: strings.NewReader(sampleLog),
		Writer: &out,
		From:   "2024-01-15T10:01:00Z",
		To:     "2024-01-15T10:03:00Z",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Scanned != 3 {
		t.Errorf("scanned: got %d, want 3", res.Scanned)
	}
	if res.Exported != 3 {
		t.Errorf("exported: got %d, want 3", res.Exported)
	}
}

func TestRun_WithIncludeFilter(t *testing.T) {
	var out bytes.Buffer
	res, err := pipeline.Run(pipeline.Config{
		Reader:   strings.NewReader(sampleLog),
		Writer:   &out,
		From:     "2024-01-15T10:00:00Z",
		To:       "2024-01-15T10:04:00Z",
		Includes: []string{"ERROR|WARN"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Matched != 2 {
		t.Errorf("matched: got %d, want 2", res.Matched)
	}
}

func TestRun_InvalidFrom(t *testing.T) {
	_, err := pipeline.Run(pipeline.Config{
		Reader: strings.NewReader(sampleLog),
		Writer: &bytes.Buffer{},
		From:   "not-a-time",
		To:     "2024-01-15T10:04:00Z",
	})
	if err == nil {
		t.Fatal("expected error for invalid from timestamp")
	}
}

func TestRun_EmptyInput(t *testing.T) {
	var out bytes.Buffer
	res, err := pipeline.Run(pipeline.Config{
		Reader: strings.NewReader(""),
		Writer: &out,
		From:   "2024-01-15T10:00:00Z",
		To:     "2024-01-15T10:04:00Z",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Scanned != 0 || res.Exported != 0 {
		t.Errorf("expected zero results, got scanned=%d exported=%d", res.Scanned, res.Exported)
	}
}
