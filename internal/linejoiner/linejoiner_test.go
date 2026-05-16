package linejoiner

import (
	"testing"
)

func TestNew_DisabledWhenNLEOne(t *testing.T) {
	j := New(1, " ")
	if j.Enabled() {
		t.Fatal("expected disabled for n=1")
	}
}

func TestNew_DisabledWhenNZero(t *testing.T) {
	j := New(0, ",")
	if j.Enabled() {
		t.Fatal("expected disabled for n=0")
	}
}

func TestNew_EnabledWhenNGtOne(t *testing.T) {
	j := New(3, "|")
	if !j.Enabled() {
		t.Fatal("expected enabled for n=3")
	}
}

func TestNew_EmptyDelimiter_DefaultsToSpace(t *testing.T) {
	j := New(2, "")
	out, ok := j.Feed("a")
	if ok {
		t.Fatalf("unexpected flush on first feed: %q", out)
	}
	out, ok = j.Feed("b")
	if !ok {
		t.Fatal("expected flush on second feed")
	}
	if out != "a b" {
		t.Fatalf("expected %q, got %q", "a b", out)
	}
}

func TestFeed_Disabled_PassesThrough(t *testing.T) {
	j := New(1, "|")
	out, ok := j.Feed("hello")
	if !ok || out != "hello" {
		t.Fatalf("expected passthrough, got %q %v", out, ok)
	}
}

func TestFeed_JoinsEveryN(t *testing.T) {
	j := New(3, "-")
	lines := []string{"a", "b", "c"}
	for i, l := range lines {
		out, ok := j.Feed(l)
		if i < 2 {
			if ok {
				t.Fatalf("line %d: unexpected flush", i)
			}
		} else {
			if !ok {
				t.Fatal("expected flush on third line")
			}
			if out != "a-b-c" {
				t.Fatalf("expected %q, got %q", "a-b-c", out)
			}
		}
	}
}

func TestFeed_MultipleGroups(t *testing.T) {
	j := New(2, "+")
	pairs := [][2]string{{"x", "y"}, {"p", "q"}}
	expected := []string{"x+y", "p+q"}
	for i, pair := range pairs {
		j.Feed(pair[0])
		out, ok := j.Feed(pair[1])
		if !ok || out != expected[i] {
			t.Fatalf("group %d: expected %q, got %q", i, expected[i], out)
		}
	}
}

func TestFlush_EmptyBuffer_ReturnsFalse(t *testing.T) {
	j := New(3, "|")
	_, ok := j.Flush()
	if ok {
		t.Fatal("expected false on empty flush")
	}
}

func TestFlush_PartialBuffer_ReturnsJoined(t *testing.T) {
	j := New(3, "|")
	j.Feed("a")
	j.Feed("b")
	out, ok := j.Flush()
	if !ok {
		t.Fatal("expected partial flush to succeed")
	}
	if out != "a|b" {
		t.Fatalf("expected %q, got %q", "a|b", out)
	}
}

func TestReset_ClearsBuffer(t *testing.T) {
	j := New(3, "|")
	j.Feed("a")
	j.Feed("b")
	j.Reset()
	_, ok := j.Flush()
	if ok {
		t.Fatal("expected empty buffer after reset")
	}
}
