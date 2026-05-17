package linetagfilter_test

import (
	"testing"

	"github.com/user/logslice/internal/linetagfilter"
)

func TestNew_EmptyTags_Disabled(t *testing.T) {
	f, err := linetagfilter.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Enabled() {
		t.Fatal("expected filter to be disabled")
	}
}

func TestNew_InvalidRegex_ReturnsError(t *testing.T) {
	_, err := linetagfilter.New([]string{"level=~[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_ValidExact_Enabled(t *testing.T) {
	f, err := linetagfilter.New([]string{"level=error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Enabled() {
		t.Fatal("expected filter to be enabled")
	}
}

func TestKeep_Disabled_AlwaysTrue(t *testing.T) {
	f, _ := linetagfilter.New(nil)
	if !f.Keep("anything here") {
		t.Fatal("disabled filter should always return true")
	}
}

func TestKeep_ExactMatch_Present(t *testing.T) {
	f, _ := linetagfilter.New([]string{"level=error"})
	line := `2024-01-01T00:00:00Z level=error msg="disk full"`
	if !f.Keep(line) {
		t.Fatal("expected line to be kept")
	}
}

func TestKeep_ExactMatch_Absent(t *testing.T) {
	f, _ := linetagfilter.New([]string{"level=error"})
	line := `2024-01-01T00:00:00Z level=info msg="ok"`
	if f.Keep(line) {
		t.Fatal("expected line to be rejected")
	}
}

func TestKeep_RegexMatch_Matches(t *testing.T) {
	f, _ := linetagfilter.New([]string{"level=~^(error|warn)$"})
	if !f.Keep("ts=1 level=warn msg=x") {
		t.Fatal("expected warn to match")
	}
	if !f.Keep("ts=1 level=error msg=x") {
		t.Fatal("expected error to match")
	}
}

func TestKeep_RegexMatch_NoMatch(t *testing.T) {
	f, _ := linetagfilter.New([]string{"level=~^(error|warn)$"})
	if f.Keep("ts=1 level=debug msg=x") {
		t.Fatal("expected debug to be rejected")
	}
}

func TestKeep_MultipleConstraints_AllMustMatch(t *testing.T) {
	f, _ := linetagfilter.New([]string{"level=error", "svc=auth"})
	if f.Keep("level=error msg=x") {
		t.Fatal("missing svc tag should be rejected")
	}
	if !f.Keep("level=error svc=auth msg=x") {
		t.Fatal("both tags present should be kept")
	}
}
