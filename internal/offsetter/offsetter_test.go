package offsetter_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/offsetter"
)

func TestNew_ZeroOffset(t *testing.T) {
	o := offsetter.New(0)
	if !o.IsZero() {
		t.Error("expected zero offset to report IsZero() == true")
	}
}

func TestNew_NonZeroOffset(t *testing.T) {
	o := offsetter.New(2 * time.Hour)
	if o.IsZero() {
		t.Error("expected non-zero offset to report IsZero() == false")
	}
}

func TestShift_PositiveOffset(t *testing.T) {
	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	o := offsetter.New(3 * time.Hour)
	got := o.Shift(base)
	want := base.Add(3 * time.Hour)
	if !got.Equal(want) {
		t.Errorf("Shift() = %v, want %v", got, want)
	}
}

func TestShift_NegativeOffset(t *testing.T) {
	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	o := offsetter.New(-1 * time.Hour)
	got := o.Shift(base)
	want := base.Add(-1 * time.Hour)
	if !got.Equal(want) {
		t.Errorf("Shift() = %v, want %v", got, want)
	}
}

func TestShiftFrom_ZeroTime_Unchanged(t *testing.T) {
	o := offsetter.New(5 * time.Hour)
	var zero time.Time
	got := o.ShiftFrom(zero)
	if !got.IsZero() {
		t.Errorf("ShiftFrom(zero) should return zero time, got %v", got)
	}
}

func TestShiftTo_ZeroTime_Unchanged(t *testing.T) {
	o := offsetter.New(5 * time.Hour)
	var zero time.Time
	got := o.ShiftTo(zero)
	if !got.IsZero() {
		t.Errorf("ShiftTo(zero) should return zero time, got %v", got)
	}
}

func TestShiftFrom_NonZero_Applied(t *testing.T) {
	base := time.Date(2024, 6, 1, 8, 30, 0, 0, time.UTC)
	o := offsetter.New(30 * time.Minute)
	got := o.ShiftFrom(base)
	want := base.Add(30 * time.Minute)
	if !got.Equal(want) {
		t.Errorf("ShiftFrom() = %v, want %v", got, want)
	}
}

func TestString_ZeroOffset(t *testing.T) {
	o := offsetter.New(0)
	if got := o.String(); got != "no offset" {
		t.Errorf("String() = %q, want %q", got, "no offset")
	}
}

func TestString_NonZeroOffset(t *testing.T) {
	o := offsetter.New(2 * time.Hour)
	if got := o.String(); got == "no offset" || got == "" {
		t.Errorf("String() = %q, expected duration string", got)
	}
}
