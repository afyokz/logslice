package cli_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/cli"
)

func withArgs(args ...string) func() {
	old := os.Args
	os.Args = append([]string{"logslice"}, args...)
	return func() { os.Args = old }
}

func TestParse_ValidArgs(t *testing.T) {
	defer withArgs(
		"--input", "app.log",
		"--from", "2024-01-01T00:00:00Z",
		"--to", "2024-01-02T00:00:00Z",
	)()

	args, err := cli.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if args.Input != "app.log" {
		t.Errorf("expected input 'app.log', got %q", args.Input)
	}
	if args.From.IsZero() || args.To.IsZero() {
		t.Error("expected From and To to be parsed")
	}
}

func TestParse_MissingInput(t *testing.T) {
	defer withArgs(
		"--from", "2024-01-01T00:00:00Z",
		"--to", "2024-01-02T00:00:00Z",
	)()

	_, err := cli.Parse()
	if err == nil {
		t.Fatal("expected error for missing --input")
	}
}

func TestParse_ToBeforeFrom(t *testing.T) {
	defer withArgs(
		"--input", "app.log",
		"--from", "2024-01-02T00:00:00Z",
		"--to", "2024-01-01T00:00:00Z",
	)()

	_, err := cli.Parse()
	if err == nil {
		t.Fatal("expected error when --to is before --from")
	}
}

func TestParse_InvalidFromFormat(t *testing.T) {
	defer withArgs(
		"--input", "app.log",
		"--from", "not-a-date",
		"--to", "2024-01-02T00:00:00Z",
	)()

	_, err := cli.Parse()
	if err == nil {
		t.Fatal("expected error for invalid --from format")
	}
}

func TestParse_IncludeExcludePatterns(t *testing.T) {
	defer withArgs(
		"--input", "app.log",
		"--from", "2024-01-01T00:00:00Z",
		"--to", "2024-01-02T00:00:00Z",
		"--include", "ERROR, FATAL",
		"--exclude", "healthcheck",
	)()

	args, err := cli.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args.IncludePatterns) != 2 {
		t.Errorf("expected 2 include patterns, got %d", len(args.IncludePatterns))
	}
	if len(args.ExcludePatterns) != 1 {
		t.Errorf("expected 1 exclude pattern, got %d", len(args.ExcludePatterns))
	}
}

func TestParse_NumberedFlag(t *testing.T) {
	defer withArgs(
		"--input", "app.log",
		"--from", "2024-01-01T00:00:00Z",
		"--to", "2024-01-02T00:00:00Z",
		"--numbered",
	)()

	args, err := cli.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !args.Numbered {
		t.Error("expected Numbered to be true")
	}
	_ = time.RFC3339 // ensure time import is used
}
