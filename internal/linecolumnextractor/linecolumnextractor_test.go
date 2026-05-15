package linecolumnextractor_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/linecolumnextractor"
)

func TestNew_EmptyColumns_Disabled(t *testing.T) {
	e := linecolumnextractor.New(" ", nil)
	if e.Enabled() {
		t.Fatal("expected extractor to be disabled when no columns supplied")
	}
}

func TestNew_WithColumns_Enabled(t *testing.T) {
	e := linecolumnextractor.New(" ", []string{"ts", "level", "msg"})
	if !e.Enabled() {
		t.Fatal("expected extractor to be enabled")
	}
}

func TestNew_EmptyDelimiter_DefaultsToSpace(t *testing.T) {
	e := linecolumnextractor.New("", []string{"a", "b"})
	result := e.Extract("hello world")
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result["a"] != "hello" || result["b"] != "world" {
		t.Fatalf("unexpected result: %v", result)
	}
}

func TestExtract_Disabled_ReturnsNil(t *testing.T) {
	e := linecolumnextractor.New(" ", nil)
	if got := e.Extract("a b c"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestExtract_InsufficientFields_ReturnsNil(t *testing.T) {
	e := linecolumnextractor.New(" ", []string{"ts", "level", "msg"})
	if got := e.Extract("only two"); got != nil {
		t.Fatalf("expected nil for insufficient fields, got %v", got)
	}
}

func TestExtract_ExactFields(t *testing.T) {
	e := linecolumnextractor.New(" ", []string{"ts", "level", "msg"})
	result := e.Extract("2024-01-01 INFO started")
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result["ts"] != "2024-01-01" {
		t.Errorf("ts: got %q", result["ts"])
	}
	if result["level"] != "INFO" {
		t.Errorf("level: got %q", result["level"])
	}
	if result["msg"] != "started" {
		t.Errorf("msg: got %q", result["msg"])
	}
}

func TestExtract_ExtraFields_CollapsedIntoLast(t *testing.T) {
	e := linecolumnextractor.New(" ", []string{"ts", "rest"})
	result := e.Extract("2024-01-01 INFO service started ok")
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result["rest"] != "INFO service started ok" {
		t.Errorf("rest: got %q", result["rest"])
	}
}

func TestExtract_CustomDelimiter(t *testing.T) {
	e := linecolumnextractor.New("|", []string{"host", "port", "status"})
	result := e.Extract("localhost|8080|200")
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result["host"] != "localhost" || result["port"] != "8080" || result["status"] != "200" {
		t.Fatalf("unexpected result: %v", result)
	}
}

func TestColumns_ReturnsCopy(t *testing.T) {
	cols := []string{"a", "b"}
	e := linecolumnextractor.New(" ", cols)
	if len(e.Columns()) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(e.Columns()))
	}
}
