package exporter_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/exporter"
)

func TestExport_Raw_Stdout(t *testing.T) {
	e, err := exporter.New(exporter.Options{Format: exporter.FormatRaw})
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer e.Close()

	lines := []string{"line one", "line two", "line three"}
	if err := e.Export(lines); err != nil {
		t.Fatalf("Export() error: %v", err)
	}
}

func TestExport_Raw_ToFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.log")
	e, err := exporter.New(exporter.Options{Format: exporter.FormatRaw, OutputPath: tmp})
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	lines := []string{"2024-01-01 foo", "2024-01-02 bar"}
	if err := e.Export(lines); err != nil {
		t.Fatalf("Export() error: %v", err)
	}
	e.Close()

	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	got := string(data)
	for _, l := range lines {
		if !strings.Contains(got, l) {
			t.Errorf("expected output to contain %q, got:\n%s", l, got)
		}
	}
}

func TestExport_Numbered_ToFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "numbered.log")
	e, err := exporter.New(exporter.Options{Format: exporter.FormatNumbered, OutputPath: tmp})
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	lines := []string{"alpha", "beta", "gamma"}
	if err := e.Export(lines); err != nil {
		t.Fatalf("Export() error: %v", err)
	}
	e.Close()

	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	got := string(data)
	if !strings.Contains(got, "1: alpha") {
		t.Errorf("expected numbered line '1: alpha' in output:\n%s", got)
	}
	if !strings.Contains(got, "3: gamma") {
		t.Errorf("expected numbered line '3: gamma' in output:\n%s", got)
	}
}

func TestExport_EmptyLines(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "empty.log")
	e, err := exporter.New(exporter.Options{OutputPath: tmp})
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer e.Close()

	if err := e.Export([]string{}); err != nil {
		t.Fatalf("Export() with empty slice should not error: %v", err)
	}
}
