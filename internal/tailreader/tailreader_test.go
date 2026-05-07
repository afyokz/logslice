package tailreader

import (
	"os"
	"testing"
	"time"
)

func writeTmp(t *testing.T, content string) *os.File {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write: %v", err)
	}
	return f
}

func TestNew_InvalidPath_ReturnsError(t *testing.T) {
	_, err := New("/nonexistent/path/file.log", 0)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestNew_DefaultPollInterval(t *testing.T) {
	f := writeTmp(t, "")
	tr, err := New(f.Name(), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer tr.Stop()
	if tr.pollInterval != 250*time.Millisecond {
		t.Errorf("expected 250ms, got %v", tr.pollInterval)
	}
}

func TestLines_ReceivesAppendedLines(t *testing.T) {
	f := writeTmp(t, "")
	tr, err := New(f.Name(), 5*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer tr.Stop()

	ch := tr.Lines()

	// Append lines after tail started
	time.Sleep(10 * time.Millisecond)
	if _, err := f.WriteString("alpha\nbeta\n"); err != nil {
		t.Fatalf("write: %v", err)
	}

	collected := make([]string, 0, 2)
	timeout := time.After(500 * time.Millisecond)
loop:
	for {
		select {
		case line, ok := <-ch:
			if !ok {
				break loop
			}
			collected = append(collected, line)
			if len(collected) == 2 {
				break loop
			}
		case <-timeout:
			break loop
		}
	}

	if len(collected) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(collected), collected)
	}
	if collected[0] != "alpha" || collected[1] != "beta" {
		t.Errorf("unexpected lines: %v", collected)
	}
}

func TestStop_ClosesChannel(t *testing.T) {
	f := writeTmp(t, "")
	tr, err := New(f.Name(), 5*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ch := tr.Lines()
	time.Sleep(20 * time.Millisecond)
	tr.Stop()

	timeout := time.After(300 * time.Millisecond)
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("expected channel to be closed")
		}
	case <-timeout:
		t.Error("channel was not closed after Stop")
	}
}

func TestIndexByte(t *testing.T) {
	if indexByte([]byte("hello\nworld"), '\n') != 5 {
		t.Error("expected index 5")
	}
	if indexByte([]byte("noeol"), '\n') != -1 {
		t.Error("expected -1")
	}
}
