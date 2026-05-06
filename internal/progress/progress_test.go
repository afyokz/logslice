package progress_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/logslice/internal/progress"
)

func TestNew_DefaultsToStderr(t *testing.T) {
	r := progress.New(0, false, nil)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}

func TestIncProcessed(t *testing.T) {
	var buf bytes.Buffer
	r := progress.New(0, false, &buf)
	r.IncProcessed()
	r.IncProcessed()
	r.IncProcessed()
	r.Print()
	if !strings.Contains(buf.String(), "processed=3") {
		t.Errorf("expected processed=3 in output, got: %s", buf.String())
	}
}

func TestIncMatched(t *testing.T) {
	var buf bytes.Buffer
	r := progress.New(0, false, &buf)
	r.IncMatched()
	r.IncMatched()
	r.Print()
	if !strings.Contains(buf.String(), "matched=2") {
		t.Errorf("expected matched=2 in output, got: %s", buf.String())
	}
}

func TestIncSkipped(t *testing.T) {
	var buf bytes.Buffer
	r := progress.New(0, false, &buf)
	r.IncSkipped()
	r.Print()
	if !strings.Contains(buf.String(), "skipped=1") {
		t.Errorf("expected skipped=1 in output, got: %s", buf.String())
	}
}

func TestPrint_Verbose_WithTotal(t *testing.T) {
	var buf bytes.Buffer
	r := progress.New(200, true, &buf)
	for i := 0; i < 100; i++ {
		r.IncProcessed()
	}
	r.IncMatched()
	r.Print()
	out := buf.String()
	if !strings.Contains(out, "50.0%") {
		t.Errorf("expected 50.0%% in verbose output, got: %s", out)
	}
	if !strings.Contains(out, "total=200") {
		t.Errorf("expected total=200 in verbose output, got: %s", out)
	}
}

func TestPrint_Verbose_NoTotal(t *testing.T) {
	var buf bytes.Buffer
	r := progress.New(0, true, &buf)
	r.IncProcessed()
	r.Print()
	out := buf.String()
	if strings.Contains(out, "total=") {
		t.Errorf("did not expect total= when total is 0, got: %s", out)
	}
}

func TestSummary(t *testing.T) {
	var buf bytes.Buffer
	r := progress.New(0, false, &buf)
	r.IncProcessed()
	r.IncProcessed()
	r.IncMatched()
	r.IncSkipped()
	s := r.Summary()
	if !strings.Contains(s, "processed=2") {
		t.Errorf("expected processed=2 in summary, got: %s", s)
	}
	if !strings.Contains(s, "matched=1") {
		t.Errorf("expected matched=1 in summary, got: %s", s)
	}
	if !strings.Contains(s, "skipped=1") {
		t.Errorf("expected skipped=1 in summary, got: %s", s)
	}
}
