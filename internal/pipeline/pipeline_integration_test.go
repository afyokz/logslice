package pipeline_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/logslice/logslice/internal/pipeline"
)

func TestRun_ExportToFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "result.log")

	_, err := pipeline.Run(pipeline.Config{
		Reader:     strings.NewReader(sampleLog),
		Writer:     nil,
		From:       "2024-01-15T10:00:00Z",
		To:         "2024-01-15T10:04:00Z",
		OutputFile: out,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if !strings.Contains(string(data), "ERROR disk full") {
		t.Errorf("output file missing expected content")
	}
}

func TestRun_NumberedOutput(t *testing.T) {
	var out bytes.Buffer
	_, err := pipeline.Run(pipeline.Config{
		Reader:  strings.NewReader(sampleLog),
		Writer:  &out,
		From:    "2024-01-15T10:00:00Z",
		To:      "2024-01-15T10:04:00Z",
		Numbered: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "1:") {
		t.Errorf("expected numbered output, got: %s", out.String())
	}
}

func TestRun_ExcludeFilter(t *testing.T) {
	var out bytes.Buffer
	res, err := pipeline.Run(pipeline.Config{
		Reader:   strings.NewReader(sampleLog),
		Writer:   &out,
		From:     "2024-01-15T10:00:00Z",
		To:       "2024-01-15T10:04:00Z",
		Excludes: []string{"DEBUG"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out.String(), "DEBUG") {
		t.Error("output should not contain excluded DEBUG lines")
	}
	if res.Matched >= res.Scanned {
		t.Errorf("expected fewer matched than scanned, got matched=%d scanned=%d", res.Matched, res.Scanned)
	}
}
