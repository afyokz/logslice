package followpipeline

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/exporter"
)

func TestRun_InvalidFile_ReturnsError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	err := Run(ctx, Config{
		FilePath:  "/no/such/file.log",
		ExportCfg: exporter.Config{Writer: os.Stdout},
	})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRun_InvalidIncludePattern_ReturnsError(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "follow-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	err = Run(ctx, Config{
		FilePath:  f.Name(),
		Include:   []string{"["},
		ExportCfg: exporter.Config{Writer: os.Stdout},
	})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestRun_ContextCancel_StopsCleanly(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "follow-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, Config{
			FilePath:     f.Name(),
			PollInterval: 5 * time.Millisecond,
			ExportCfg:    exporter.Config{Writer: os.Stdout},
		})
	}()

	time.Sleep(30 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != nil && err != context.Canceled {
			t.Errorf("unexpected error: %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("Run did not stop after context cancel")
	}
}

func TestRun_FiltersAndExports(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "follow-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var buf strings.Builder
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, Config{
			FilePath:     f.Name(),
			PollInterval: 5 * time.Millisecond,
			Include:      []string{"ERROR"},
			ExportCfg:    exporter.Config{Writer: &buf},
		})
	}()

	time.Sleep(20 * time.Millisecond)
	f.WriteString("INFO  all good\n")
	f.WriteString("ERROR something failed\n")
	time.Sleep(50 * time.Millisecond)
	cancel()
	<-done

	if !strings.Contains(buf.String(), "ERROR something failed") {
		t.Errorf("expected ERROR line in output, got: %q", buf.String())
	}
	if strings.Contains(buf.String(), "INFO") {
		t.Errorf("INFO line should have been filtered out, got: %q", buf.String())
	}
}
