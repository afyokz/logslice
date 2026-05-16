package linetemplate_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/linetemplate"
)

func TestNew_EmptyTemplate_Disabled(t *testing.T) {
	f, err := linetemplate.New("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Enabled() {
		t.Fatal("expected formatter to be disabled for empty template")
	}
}

func TestNew_InvalidTemplate_ReturnsError(t *testing.T) {
	_, err := linetemplate.New("{{.Foo")
	if err == nil {
		t.Fatal("expected error for invalid template syntax")
	}
}

func TestNew_ValidTemplate_Enabled(t *testing.T) {
	f, err := linetemplate.New("{{.level}}: {{.msg}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Enabled() {
		t.Fatal("expected formatter to be enabled")
	}
}

func TestFormat_Disabled_ReturnsOriginal(t *testing.T) {
	f, _ := linetemplate.New("")
	got := f.Format("original line", map[string]string{"level": "INFO"})
	if got != "original line" {
		t.Errorf("expected original line, got %q", got)
	}
}

func TestFormat_AppliesTemplate(t *testing.T) {
	f, err := linetemplate.New("[{{.level}}] {{.msg}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fields := map[string]string{"level": "WARN", "msg": "disk full"}
	got := f.Format("raw line", fields)
	want := "[WARN] disk full"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestFormat_MissingKey_ReturnsEmptyForField(t *testing.T) {
	f, err := linetemplate.New("{{.level}}: {{.msg}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fields := map[string]string{"level": "ERROR"}
	got := f.Format("raw", fields)
	want := "ERROR: "
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestFormat_NilFields_ReturnsRendered(t *testing.T) {
	f, err := linetemplate.New("static output")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := f.Format("raw", nil)
	if got != "static output" {
		t.Errorf("expected %q, got %q", "static output", got)
	}
}
