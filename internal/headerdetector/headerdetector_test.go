package headerdetector

import (
	"testing"
)

func TestNew_DefaultPatterns(t *testing.T) {
	d, err := New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !d.enabled {
		t.Fatal("expected detector to be enabled")
	}
	if len(d.patterns) != len(defaultPatterns) {
		t.Fatalf("expected %d patterns, got %d", len(defaultPatterns), len(d.patterns))
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := New([]string{`(?P<bad`})
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestIsHeader_ISO8601(t *testing.T) {
	d, _ := New(nil)
	line := "2024-05-01T12:34:56 INFO starting server"
	if !d.IsHeader(line) {
		t.Errorf("expected ISO-8601 line to be a header")
	}
}

func TestIsHeader_Syslog(t *testing.T) {
	d, _ := New(nil)
	line := "May  1 12:34:56 myhostname myapp[123]: message"
	if !d.IsHeader(line) {
		t.Errorf("expected syslog line to be a header")
	}
}

func TestIsContinuation_StackTrace(t *testing.T) {
	d, _ := New(nil)
	line := "\tat com.example.App.main(App.java:42)"
	if !d.IsContinuation(line) {
		t.Errorf("expected indented stack-trace line to be a continuation")
	}
}

func TestIsContinuation_PlainText(t *testing.T) {
	d, _ := New(nil)
	if !d.IsContinuation("just some random text") {
		t.Errorf("expected plain text to be a continuation")
	}
}

func TestIsHeader_CustomPattern(t *testing.T) {
	d, err := New([]string{`^ERROR`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !d.IsHeader("ERROR something went wrong") {
		t.Error("expected custom-pattern line to be a header")
	}
	if d.IsHeader("INFO all good") {
		t.Error("expected non-matching line NOT to be a header")
	}
}

func TestStripTimestamp_RemovesPrefix(t *testing.T) {
	d, _ := New(nil)
	line := "2024-05-01T12:34:56 INFO server started"
	want := "INFO server started"
	got := d.StripTimestamp(line)
	if got != want {
		t.Errorf("StripTimestamp: got %q, want %q", got, want)
	}
}

func TestStripTimestamp_NoMatch_ReturnsOriginal(t *testing.T) {
	d, _ := New(nil)
	line := "no timestamp here"
	if got := d.StripTimestamp(line); got != line {
		t.Errorf("expected original line, got %q", got)
	}
}

func TestDisabled_AllLinesAreHeaders(t *testing.T) {
	d := &Detector{enabled: false}
	if !d.IsHeader("anything") {
		t.Error("disabled detector should treat every line as a header")
	}
	if d.IsContinuation("anything") {
		t.Error("disabled detector should never report a continuation")
	}
}
