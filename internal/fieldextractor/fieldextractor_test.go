package fieldextractor_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/fieldextractor"
)

func TestNew_EmptyPattern_Disabled(t *testing.T) {
	e, err := fieldextractor.New("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Enabled() {
		t.Error("expected extractor to be disabled for empty pattern")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := fieldextractor.New(`(?P<bad[`)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestNew_ValidPattern_Enabled(t *testing.T) {
	e, err := fieldextractor.New(`(?P<level>\w+)\s+(?P<msg>.+)`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !e.Enabled() {
		t.Error("expected extractor to be enabled")
	}
}

func TestFields_ReturnsNamedGroups(t *testing.T) {
	e, _ := fieldextractor.New(`(?P<level>\w+)\s+(?P<msg>.+)`)
	fields := e.Fields()
	if len(fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(fields))
	}
	if fields[0] != "level" || fields[1] != "msg" {
		t.Errorf("unexpected fields: %v", fields)
	}
}

func TestExtract_Disabled_ReturnsNil(t *testing.T) {
	e, _ := fieldextractor.New("")
	result := e.Extract("INFO something happened")
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestExtract_NoMatch_ReturnsNil(t *testing.T) {
	e, _ := fieldextractor.New(`^\[(?P<ts>[\d\-T:]+)\]`)
	result := e.Extract("no timestamp here")
	if result != nil {
		t.Errorf("expected nil for non-matching line, got %v", result)
	}
}

func TestExtract_Match_ReturnsFields(t *testing.T) {
	e, _ := fieldextractor.New(`(?P<level>\w+)\s+(?P<msg>.+)`)
	result := e.Extract("ERROR disk full")
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result["level"] != "ERROR" {
		t.Errorf("expected level=ERROR, got %q", result["level"])
	}
	if result["msg"] != "disk full" {
		t.Errorf("expected msg='disk full', got %q", result["msg"])
	}
}

func TestExtract_PartialMatch_CapturesAvailableFields(t *testing.T) {
	e, _ := fieldextractor.New(`(?P<ip>\d+\.\d+\.\d+\.\d+)\s+(?P<user>\S+)`)
	result := e.Extract("192.168.1.1 alice GET /index")
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result["ip"] != "192.168.1.1" {
		t.Errorf("unexpected ip: %q", result["ip"])
	}
	if result["user"] != "alice" {
		t.Errorf("unexpected user: %q", result["user"])
	}
}
