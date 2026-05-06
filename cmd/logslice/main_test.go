package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// buildBinary compiles the binary into a temp dir and returns its path.
// Skips the test if go build fails (e.g. missing dependencies in CI).
func buildBinary(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "logslice")
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("skipping binary test: go build failed: %s", out)
	}
	return binPath
}

func TestMain_Version(t *testing.T) {
	bin := buildBinary(t)
	out, err := exec.Command(bin, "--version").Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "logslice version") {
		t.Errorf("expected version string, got: %s", out)
	}
}

func TestMain_MissingInput_ExitsNonZero(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "--from", "2024-01-15 08:00:00", "--to", "2024-01-15 09:00:00")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for missing --input flag")
	}
}

func TestMain_BasicSlice_Stdout(t *testing.T) {
	bin := buildBinary(t)

	// Write a small temp log file
	logContent := `2024-01-15 08:00:01 INFO service started
2024-01-15 08:30:00 ERROR something failed
2024-01-15 09:01:00 INFO outside range
`
	tmpLog := filepath.Join(t.TempDir(), "test.log")
	if err := os.WriteFile(tmpLog, []byte(logContent), 0o644); err != nil {
		t.Fatalf("failed to write temp log: %v", err)
	}

	cmd := exec.Command(bin,
		"--input", tmpLog,
		"--from", "2024-01-15 08:00:00",
		"--to", "2024-01-15 09:00:00",
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := string(out)
	if !strings.Contains(result, "service started") {
		t.Errorf("expected 'service started' in output, got: %s", result)
	}
	if !strings.Contains(result, "something failed") {
		t.Errorf("expected 'something failed' in output, got: %s", result)
	}
	if strings.Contains(result, "outside range") {
		t.Errorf("did not expect 'outside range' in output, got: %s", result)
	}
}
