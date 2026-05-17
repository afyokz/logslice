package linescaler_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/linescaler"
)

func TestNew_DisabledWhenFactorZero(t *testing.T) {
	s := linescaler.New(0, 0, " ")
	if s.Enabled() {
		t.Fatal("expected disabled when factor is 0")
	}
}

func TestNew_DisabledWhenFactorOne(t *testing.T) {
	s := linescaler.New(1, 0, " ")
	if s.Enabled() {
		t.Fatal("expected disabled when factor is 1")
	}
}

func TestNew_EnabledWhenFactorOther(t *testing.T) {
	s := linescaler.New(2.5, 0, " ")
	if !s.Enabled() {
		t.Fatal("expected enabled when factor != 0 and != 1")
	}
}

func TestNew_DefaultDelimiter(t *testing.T) {
	s := linescaler.New(2, 0, "")
	if !s.Enabled() {
		t.Fatal("expected enabled")
	}
	out, err := s.Scale("10 hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "20 hello" {
		t.Fatalf("got %q, want %q", out, "20 hello")
	}
}

func TestScale_Disabled_ReturnsOriginal(t *testing.T) {
	s := linescaler.New(1, 0, " ")
	line := "42 some log line"
	out, err := s.Scale(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != line {
		t.Fatalf("got %q, want %q", out, line)
	}
}

func TestScale_FirstField(t *testing.T) {
	s := linescaler.New(10, 0, " ")
	out, err := s.Scale("5 INFO something happened")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "50 INFO something happened" {
		t.Fatalf("got %q", out)
	}
}

func TestScale_MiddleField(t *testing.T) {
	s := linescaler.New(0.5, 2, ",")
	out, err := s.Scale("ts,INFO,100,msg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "ts,INFO,50,msg" {
		t.Fatalf("got %q", out)
	}
}

func TestScale_FieldOutOfRange_ReturnsError(t *testing.T) {
	s := linescaler.New(2, 5, " ")
	original := "only two fields"
	out, err := s.Scale(original)
	if err == nil {
		t.Fatal("expected error for out-of-range field")
	}
	if out != original {
		t.Fatalf("expected original line on error, got %q", out)
	}
}

func TestScale_NonNumericField_ReturnsError(t *testing.T) {
	s := linescaler.New(2, 1, " ")
	original := "ts notanumber msg"
	out, err := s.Scale(original)
	if err == nil {
		t.Fatal("expected error for non-numeric field")
	}
	if out != original {
		t.Fatalf("expected original line on error, got %q", out)
	}
}

func TestFactor_ReturnsConfiguredValue(t *testing.T) {
	s := linescaler.New(3.14, 0, " ")
	if s.Factor() != 3.14 {
		t.Fatalf("got %v, want 3.14", s.Factor())
	}
}
