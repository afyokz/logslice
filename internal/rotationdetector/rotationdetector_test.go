package rotationdetector_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logslice/internal/rotationdetector"
)

func TestNew_EmptyPath_Disabled(t *testing.T) {
	d, err := rotationdetector.New("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Enabled() {
		t.Error("expected detector to be disabled for empty path")
	}
}

func TestNew_NonExistentPath_ReturnsError(t *testing.T) {
	_, err := rotationdetector.New("/nonexistent/path/to/log.txt")
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

func TestNew_ValidFile_Enabled(t *testing.T) {
	f := tmpFile(t, "hello\n")
	d, err := rotationdetector.New(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !d.Enabled() {
		t.Error("expected detector to be enabled")
	}
}

func TestRotated_NoChange_ReturnsFalse(t *testing.T) {
	f := tmpFile(t, "line1\n")
	d, err := rotationdetector.New(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rotated, err := d.Rotated()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rotated {
		t.Error("expected Rotated=false when file unchanged")
	}
}

func TestRotated_Truncated_ReturnsTrue(t *testing.T) {
	f := tmpFile(t, "line1\nline2\nline3\n")
	d, err := rotationdetector.New(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Truncate the file to simulate log rotation.
	if err := os.WriteFile(f, []byte(""), 0644); err != nil {
		t.Fatalf("failed to truncate file: %v", err)
	}
	rotated, err := d.Rotated()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rotated {
		t.Error("expected Rotated=true after truncation")
	}
}

func TestRotated_Deleted_ReturnsTrue(t *testing.T) {
	f := tmpFile(t, "data\n")
	d, err := rotationdetector.New(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	os.Remove(f)
	rotated, err := d.Rotated()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rotated {
		t.Error("expected Rotated=true after file deletion")
	}
}

func TestReset_ClearsRotationState(t *testing.T) {
	f := tmpFile(t, "line1\nline2\n")
	d, err := rotationdetector.New(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Truncate, confirm rotated, then reset.
	os.WriteFile(f, []byte("new\n"), 0644)
	if err := d.Reset(); err != nil {
		t.Fatalf("Reset error: %v", err)
	}
	rotated, err := d.Rotated()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rotated {
		t.Error("expected Rotated=false after Reset")
	}
}

func tmpFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "test.log")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	return p
}
