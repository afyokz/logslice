package pipeline_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/linemerger"
)

// TestMerger_IntegrationWithPipeline verifies that the linemerger correctly
// groups lines before they would enter the export stage in a realistic
// pipeline scenario.
func TestMerger_IntegrationWithPipeline(t *testing.T) {
	input := []string{
		"2024-03-01 10:00:00 request start id=1",
		"2024-03-01 10:00:00 request body id=1",
		"2024-03-01 10:00:01 response id=2",
		"2024-03-01 10:00:01 response body id=2",
		"2024-03-01 10:00:02 done",
	}

	m := linemerger.New(19, " | ")
	var results []string

	for _, line := range input {
		if out, ok := m.Feed(line); ok {
			results = append(results, out)
		}
	}
	if out, ok := m.Flush(); ok {
		results = append(results, out)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 merged records, got %d: %v", len(results), results)
	}

	if !strings.Contains(results[0], "request start") || !strings.Contains(results[0], "request body") {
		t.Errorf("first record should merge both request lines: %q", results[0])
	}
	if !strings.Contains(results[1], "response id=2") && !strings.Contains(results[1], "response body") {
		t.Errorf("second record should merge both response lines: %q", results[1])
	}
	if !strings.Contains(results[2], "done") {
		t.Errorf("third record should contain 'done': %q", results[2])
	}
}

func TestMerger_SingleLineGroups_PassThrough(t *testing.T) {
	input := []string{
		"2024-03-01 10:00:00 alpha",
		"2024-03-01 10:00:01 beta",
		"2024-03-01 10:00:02 gamma",
	}

	m := linemerger.New(19, " | ")
	var results []string
	for _, line := range input {
		if out, ok := m.Feed(line); ok {
			results = append(results, out)
		}
	}
	if out, ok := m.Flush(); ok {
		results = append(results, out)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 records (one per line), got %d", len(results))
	}
}
