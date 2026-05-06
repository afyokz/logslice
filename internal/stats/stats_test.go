package stats_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/stats"
)

func TestCounter_Increments(t *testing.T) {
	var c stats.Counter
	c.IncRead()
	c.IncRead()
	c.IncMatched()
	c.IncSkipped()
	c.IncExported()

	if c.LinesRead != 2 {
		t.Errorf("LinesRead: want 2, got %d", c.LinesRead)
	}
	if c.LinesMatched != 1 {
		t.Errorf("LinesMatched: want 1, got %d", c.LinesMatched)
	}
	if c.LinesSkipped != 1 {
		t.Errorf("LinesSkipped: want 1, got %d", c.LinesSkipped)
	}
	if c.LinesExported != 1 {
		t.Errorf("LinesExported: want 1, got %d", c.LinesExported)
	}
}

func TestCounter_Add(t *testing.T) {
	a := stats.Counter{LinesRead: 5, LinesMatched: 3, LinesSkipped: 1, LinesExported: 3}
	b := stats.Counter{LinesRead: 2, LinesMatched: 1, LinesSkipped: 0, LinesExported: 1}
	a.Add(b)

	if a.LinesRead != 7 {
		t.Errorf("LinesRead after Add: want 7, got %d", a.LinesRead)
	}
	if a.LinesMatched != 4 {
		t.Errorf("LinesMatched after Add: want 4, got %d", a.LinesMatched)
	}
	if a.LinesExported != 4 {
		t.Errorf("LinesExported after Add: want 4, got %d", a.LinesExported)
	}
}

func TestCounter_Print(t *testing.T) {
	c := stats.Counter{LinesRead: 10, LinesMatched: 7, LinesSkipped: 3, LinesExported: 7}
	var buf bytes.Buffer
	c.Print(&buf)
	out := buf.String()

	for _, want := range []string{"10", "7", "3", "exported"} {
		if !strings.Contains(out, want) {
			t.Errorf("Print output missing %q; got:\n%s", want, out)
		}
	}
}

func TestCounter_ZeroValue(t *testing.T) {
	var c stats.Counter
	var buf bytes.Buffer
	c.Print(&buf)
	out := buf.String()
	if !strings.Contains(out, "0") {
		t.Errorf("expected zeroes in output, got: %s", out)
	}
}
