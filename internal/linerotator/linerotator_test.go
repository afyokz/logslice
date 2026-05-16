package linerotator

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestNew_DisabledWhenMaxLinesZero(t *testing.T) {
	var buf bytes.Buffer
	r, err := New(&buf, nil, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.enabled {
		t.Fatal("expected rotator to be disabled")
	}
}

func TestNew_DisabledWhenMaxLinesNegative(t *testing.T) {
	var buf bytes.Buffer
	r, err := New(&buf, nil, -5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.enabled {
		t.Fatal("expected rotator to be disabled")
	}
}

func TestNew_FactoryErrorPropagated(t *testing.T) {
	fail := func(_ int) (interface{ Write([]byte) (int, error) }, error) {
		return nil, errors.New("factory fail")
	}
	_, err := New(nil, func(seq int) (interface{ Write([]byte) (int, error) }, error) {
		return fail(seq)
	}, 3)
	// adapt to io.Writer signature via helper below
	_, err2 := New(nil, func(seq int) (writerIface, error) {
		return nil, errors.New("factory fail")
	}, 3)
	if err2 == nil && err == nil {
		t.Fatal("expected error from factory")
	}
}

type writerIface = interface{ Write([]byte) (int, error) }

func TestNew_EnabledCallsFactoryOnce(t *testing.T) {
	calls := 0
	bufs := make([]*bytes.Buffer, 0)
	factory := func(_ int) (*bytes.Buffer, error) {
		calls++
		b := &bytes.Buffer{}
		bufs = append(bufs, b)
		return b, nil
	}
	_, err := New(nil, func(seq int) (writerIface, error) { return factory(seq) }, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 factory call, got %d", calls)
	}
}

func TestWriteLine_Disabled_AllToFallback(t *testing.T) {
	var buf bytes.Buffer
	r, _ := New(&buf, nil, 0)
	for _, line := range []string{"a", "b", "c"} {
		if err := r.WriteLine(line); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if got := buf.String(); got != "a\nb\nc\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestWriteLine_RotatesAfterMaxLines(t *testing.T) {
	bufs := []*bytes.Buffer{}
	factory := func(_ int) (writerIface, error) {
		b := &bytes.Buffer{}
		bufs = append(bufs, b)
		return b, nil
	}
	r, err := New(nil, factory, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := []string{"line1", "line2", "line3", "line4", "line5"}
	for _, l := range lines {
		if err := r.WriteLine(l); err != nil {
			t.Fatalf("write error: %v", err)
		}
	}
	if r.Written() != 5 {
		t.Fatalf("expected 5 written, got %d", r.Written())
	}
	if r.Seq() != 2 {
		t.Fatalf("expected seq 2, got %d", r.Seq())
	}
	if len(bufs) != 3 {
		t.Fatalf("expected 3 buffers, got %d", len(bufs))
	}
	expected := []string{"line1\nline2\n", "line3\nline4\n", "line5\n"}
	for i, b := range bufs {
		if got := b.String(); got != expected[i] {
			t.Errorf("buf[%d]: got %q, want %q", i, got, expected[i])
		}
	}
}

func TestWriteLine_RotateFactoryError_Propagated(t *testing.T) {
	call := 0
	factory := func(_ int) (writerIface, error) {
		call++
		if call > 1 {
			return nil, errors.New("rotate fail")
		}
		return &bytes.Buffer{}, nil
	}
	r, _ := New(nil, factory, 1)
	_ = r.WriteLine("first")
	err := r.WriteLine("second")
	if err == nil || !strings.Contains(err.Error(), "rotate fail") {
		t.Fatalf("expected rotate error, got %v", err)
	}
}
