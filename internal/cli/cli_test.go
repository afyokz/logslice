package cli

import (
	"os"
	"testing"
	"time"
)

func withArgs(args []string, fn func()) {
	old := os.Args
	os.Args = append([]string{"logslice"}, args...)
	defer func() { os.Args = old }()
	fn()
}

func TestParse_ValidArgs(t *testing.T) {
	withArgs([]string{
		"--input", "app.log",
		"--from", "2024-01-15T08:00:00Z",
		"--to", "2024-01-15T09:00:00Z",
		"--numbered",
	}, func() {
		cfg, err := Parse()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.InputFile != "app.log" {
			t.Errorf("expected InputFile=app.log, got %q", cfg.InputFile)
		}
		if !cfg.Numbered {
			t.Error("expected Numbered=true")
		}
		expectedFrom, _ := time.Parse(time.RFC3339, "2024-01-15T08:00:00Z")
		if !cfg.From.Equal(expectedFrom) {
			t.Errorf("From mismatch: got %v", cfg.From)
		}
	})
}

func TestParse_MissingInput(t *testing.T) {
	withArgs([]string{
		"--from", "2024-01-15T08:00:00Z",
		"--to", "2024-01-15T09:00:00Z",
	}, func() {
		_, err := Parse()
		if err == nil {
			t.Fatal("expected error for missing --input")
		}
	})
}

func TestParse_ToBeforeFrom(t *testing.T) {
	withArgs([]string{
		"--input", "app.log",
		"--from", "2024-01-15T09:00:00Z",
		"--to", "2024-01-15T08:00:00Z",
	}, func() {
		_, err := Parse()
		if err == nil {
			t.Fatal("expected error when --to is before --from")
		}
	})
}

func TestParse_InvalidFromFormat(t *testing.T) {
	withArgs([]string{
		"--input", "app.log",
		"--from", "not-a-date",
		"--to", "2024-01-15T09:00:00Z",
	}, func() {
		_, err := Parse()
		if err == nil {
			t.Fatal("expected error for invalid --from")
		}
	})
}

func TestParse_CustomFormat(t *testing.T) {
	withArgs([]string{
		"--input", "app.log",
		"--from", "2024/01/15 08:00:00",
		"--to", "2024/01/15 09:00:00",
		"--format", "2006/01/02 15:04:05",
	}, func() {
		cfg, err := Parse()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Format != "2006/01/02 15:04:05" {
			t.Errorf("expected custom format, got %q", cfg.Format)
		}
	})
}
