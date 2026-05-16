package linecounter_test

import (
	"testing"

	"github.com/logslice/logslice/internal/linecounter"
)

func TestNew_DisabledCounter(t *testing.T) {
	c := linecounter.New(false)
	if c.Enabled() {
		t.Fatal("expected counter to be disabled")
	}
}

func TestNew_EnabledCounter(t *testing.T) {
	c := linecounter.New(true)
	if !c.Enabled() {
		t.Fatal("expected counter to be enabled")
	}
}

func TestInc_Disabled_NoOp(t *testing.T) {
	c := linecounter.New(false)
	c.Inc("errors")
	if got := c.Get("errors"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestInc_Increments(t *testing.T) {
	c := linecounter.New(true)
	c.Inc("errors")
	c.Inc("errors")
	c.Inc("warnings")
	if got := c.Get("errors"); got != 2 {
		t.Fatalf("errors: expected 2, got %d", got)
	}
	if got := c.Get("warnings"); got != 1 {
		t.Fatalf("warnings: expected 1, got %d", got)
	}
}

func TestAdd_Disabled_NoOp(t *testing.T) {
	c := linecounter.New(false)
	c.Add("info", 5)
	if got := c.Get("info"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAdd_NegativeOrZero_NoOp(t *testing.T) {
	c := linecounter.New(true)
	c.Add("info", 0)
	c.Add("info", -3)
	if got := c.Get("info"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAdd_AccumulatesValue(t *testing.T) {
	c := linecounter.New(true)
	c.Add("debug", 10)
	c.Add("debug", 5)
	if got := c.Get("debug"); got != 15 {
		t.Fatalf("expected 15, got %d", got)
	}
}

func TestGet_UnknownLabel_ReturnsZero(t *testing.T) {
	c := linecounter.New(true)
	if got := c.Get("nonexistent"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestTotal_SumsAllCounts(t *testing.T) {
	c := linecounter.New(true)
	c.Add("a", 3)
	c.Add("b", 7)
	c.Inc("c")
	if got := c.Total(); got != 11 {
		t.Fatalf("expected 11, got %d", got)
	}
}

func TestReset_ClearsCounts(t *testing.T) {
	c := linecounter.New(true)
	c.Inc("x")
	c.Reset()
	if got := c.Get("x"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
	if got := c.Total(); got != 0 {
		t.Fatalf("expected total 0 after reset, got %d", got)
	}
}

func TestLabels_ReturnsRecordedLabels(t *testing.T) {
	c := linecounter.New(true)
	c.Inc("alpha")
	c.Inc("beta")
	labels := c.Labels()
	if len(labels) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(labels))
	}
}
