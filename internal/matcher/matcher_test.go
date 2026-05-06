package matcher_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/matcher"
)

func TestNew_InvalidIncludePattern(t *testing.T) {
	_, err := matcher.New(matcher.Config{
		IncludePatterns: []string{"["},
	})
	if err == nil {
		t.Fatal("expected error for invalid include pattern, got nil")
	}
}

func TestNew_InvalidExcludePattern(t *testing.T) {
	_, err := matcher.New(matcher.Config{
		ExcludePatterns: []string{"("},
	})
	if err == nil {
		t.Fatal("expected error for invalid exclude pattern, got nil")
	}
}

func TestMatch_NoPatterns(t *testing.T) {
	m, err := matcher.New(matcher.Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Match("any line") {
		t.Error("expected empty matcher to pass all lines")
	}
}

func TestMatch_IncludePattern(t *testing.T) {
	m, _ := matcher.New(matcher.Config{
		IncludePatterns: []string{"ERROR"},
	})

	if !m.Match("2024-01-01 ERROR something failed") {
		t.Error("expected line with ERROR to match")
	}
	if m.Match("2024-01-01 INFO all good") {
		t.Error("expected line without ERROR to not match")
	}
}

func TestMatch_ExcludePattern(t *testing.T) {
	m, _ := matcher.New(matcher.Config{
		ExcludePatterns: []string{"DEBUG"},
	})

	if m.Match("2024-01-01 DEBUG verbose output") {
		t.Error("expected DEBUG line to be excluded")
	}
	if !m.Match("2024-01-01 INFO startup complete") {
		t.Error("expected non-DEBUG line to pass")
	}
}

func TestMatch_ExcludeTakesPriority(t *testing.T) {
	m, _ := matcher.New(matcher.Config{
		IncludePatterns: []string{"ERROR"},
		ExcludePatterns: []string{"noisy"},
	})

	if m.Match("ERROR noisy component") {
		t.Error("exclude should take priority over include")
	}
	if !m.Match("ERROR real failure") {
		t.Error("non-excluded ERROR line should match")
	}
}

func TestMatchAll(t *testing.T) {
	m, _ := matcher.New(matcher.Config{
		IncludePatterns: []string{"ERROR", "db"},
	})

	if !m.MatchAll("ERROR in db connection") {
		t.Error("expected line matching all patterns to pass MatchAll")
	}
	if m.MatchAll("ERROR in network") {
		t.Error("expected line missing 'db' to fail MatchAll")
	}
}

func TestContainsAny(t *testing.T) {
	if !matcher.ContainsAny("fatal error occurred", []string{"warn", "fatal"}) {
		t.Error("expected ContainsAny to return true")
	}
	if matcher.ContainsAny("info: all good", []string{"warn", "fatal"}) {
		t.Error("expected ContainsAny to return false")
	}
}
