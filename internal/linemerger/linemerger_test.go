package linemerger

import (
	"testing"
)

func TestNew_DisabledWhenZero(t *testing.T) {
	m := New(0, " ")
	if m.Enabled() {
		t.Fatal("expected disabled")
	}
}

func TestNew_DisabledWhenNegative(t *testing.T) {
	m := New(-1, " ")
	if m.Enabled() {
		t.Fatal("expected disabled")
	}
}

func TestNew_EnabledWhenPositive(t *testing.T) {
	m := New(10, " ")
	if !m.Enabled() {
		t.Fatal("expected enabled")
	}
}

func TestNew_DefaultDelimiter(t *testing.T) {
	m := New(5, "")
	if m.delimiter != " " {
		t.Fatalf("expected default delimiter ' ', got %q", m.delimiter)
	}
}

func TestFeed_Disabled_PassesThrough(t *testing.T) {
	m := New(0, " ")
	line := "2024-01-01 some message"
	out, ok := m.Feed(line)
	if !ok || out != line {
		t.Fatalf("expected passthrough, got %q %v", out, ok)
	}
}

func TestFeed_SameKey_Buffered(t *testing.T) {
	m := New(10, "|")
	_, ok := m.Feed("2024-01-01 first")
	if ok {
		t.Fatal("first line of group should be buffered, not emitted")
	}
	_, ok = m.Feed("2024-01-01 second")
	if ok {
		t.Fatal("second line with same key should still be buffered")
	}
}

func TestFeed_NewKey_FlushesGroup(t *testing.T) {
	m := New(10, "|")
	m.Feed("2024-01-01 first")
	m.Feed("2024-01-01 second")
	out, ok := m.Feed("2024-01-02 third")
	if !ok {
		t.Fatal("new key should flush previous group")
	}
	expected := "2024-01-01 first|2024-01-01 second"
	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func TestFlush_ReturnsMergedBuffer(t *testing.T) {
	m := New(10, " ")
	m.Feed("2024-01-01 alpha")
	m.Feed("2024-01-01 beta")
	out, ok := m.Flush()
	if !ok {
		t.Fatal("flush should return buffered content")
	}
	if out != "2024-01-01 alpha 2024-01-01 beta" {
		t.Fatalf("unexpected merged output: %q", out)
	}
}

func TestFlush_Empty_ReturnsFalse(t *testing.T) {
	m := New(10, " ")
	_, ok := m.Flush()
	if ok {
		t.Fatal("flush on empty buffer should return false")
	}
}

func TestReset_ClearsBuffer(t *testing.T) {
	m := New(10, " ")
	m.Feed("2024-01-01 something")
	m.Reset()
	_, ok := m.Flush()
	if ok {
		t.Fatal("after reset flush should return nothing")
	}
}

func TestFeed_ShortLine_UsesFullLine(t *testing.T) {
	m := New(20, "|")
	m.Feed("short")
	out, ok := m.Feed("short")
	if ok {
		t.Fatalf("same short-line key should buffer, got %q", out)
	}
	out, ok = m.Flush()
	if !ok || out != "short|short" {
		t.Fatalf("unexpected: %q %v", out, ok)
	}
}
