package linesplitter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/linesplitter"
)

func TestNew_EmptyFields_Disabled(t *testing.T) {
	s := linesplitter.New(" ", nil)
	if s.Enabled() {
		t.Fatal("expected splitter to be disabled when no fields provided")
	}
}

func TestNew_WithFields_Enabled(t *testing.T) {
	s := linesplitter.New(" ", []string{"level", "msg"})
	if !s.Enabled() {
		t.Fatal("expected splitter to be enabled")
	}
}

func TestNew_EmptyDelimiter_DefaultsToSpace(t *testing.T) {
	s := linesplitter.New("", []string{"a", "b"})
	result := s.Extract("hello world")
	if result["a"] != "hello" || result["b"] != "world" {
		t.Fatalf("unexpected result: %v", result)
	}
}

func TestExtract_Disabled_ReturnsNil(t *testing.T) {
	s := linesplitter.New(" ", nil)
	if got := s.Extract("some line"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestExtract_ExactFields(t *testing.T) {
	s := linesplitter.New(" ", []string{"level", "ts", "msg"})
	result := s.Extract("INFO 2024-01-02 hello world")
	if result["level"] != "INFO" {
		t.Errorf("level: got %q", result["level"])
	}
	if result["ts"] != "2024-01-02" {
		t.Errorf("ts: got %q", result["ts"])
	}
	// SplitN with n=3 puts the remainder in the last field
	if result["msg"] != "hello world" {
		t.Errorf("msg: got %q", result["msg"])
	}
}

func TestExtract_FewerPartsThanFields_EmptyStrings(t *testing.T) {
	s := linesplitter.New("|", []string{"a", "b", "c"})
	result := s.Extract("foo|bar")
	if result["a"] != "foo" {
		t.Errorf("a: got %q", result["a"])
	}
	if result["b"] != "bar" {
		t.Errorf("b: got %q", result["b"])
	}
	if result["c"] != "" {
		t.Errorf("c: expected empty, got %q", result["c"])
	}
}

func TestExtract_CustomDelimiter(t *testing.T) {
	s := linesplitter.New(",", []string{"host", "port"})
	result := s.Extract("localhost,8080")
	if result["host"] != "localhost" || result["port"] != "8080" {
		t.Fatalf("unexpected result: %v", result)
	}
}

func TestFields_ReturnsConfiguredNames(t *testing.T) {
	names := []string{"x", "y", "z"}
	s := linesplitter.New(" ", names)
	got := s.Fields()
	if len(got) != len(names) {
		t.Fatalf("expected %d fields, got %d", len(names), len(got))
	}
	for i, n := range names {
		if got[i] != n {
			t.Errorf("field[%d]: want %q got %q", i, n, got[i])
		}
	}
}
