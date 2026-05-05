package timeparser

import (
	"testing"
	"time"
)

func TestParseTimestamp(t *testing.T) {
	cases := []struct {
		input    string
		wantYear int
		wantErr  bool
	}{
		{"2024-03-15T10:22:33Z", 2024, false},
		{"2024-03-15T10:22:33.123456789Z", 2024, false},
		{"2024-03-15T10:22:33", 2024, false},
		{"2024-03-15 10:22:33", 2024, false},
		{"2024-03-15 10:22:33.000", 2024, false},
		{"15/Mar/2024:10:22:33 +0000", 2024, false},
		{"not-a-timestamp", 0, true},
		{"", 0, true},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, _, err := ParseTimestamp(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error for input %q, got none", tc.input)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for input %q: %v", tc.input, err)
				return
			}
			if got.Year() != tc.wantYear {
				t.Errorf("year mismatch: got %d, want %d", got.Year(), tc.wantYear)
			}
		})
	}
}

func TestInRange(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 23, 59, 59, 0, time.UTC)

	cases := []struct {
		name string
		t    time.Time
		want bool
	}{
		{"before range", start.Add(-time.Second), false},
		{"at start", start, true},
		{"in middle", start.Add(12 * time.Hour), true},
		{"at end", end, true},
		{"after range", end.Add(time.Second), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := InRange(tc.t, start, end)
			if got != tc.want {
				t.Errorf("InRange(%v) = %v, want %v", tc.t, got, tc.want)
			}
		})
	}
}

func TestParseWithFormat(t *testing.T) {
	_, err := ParseWithFormat("2024-03-15 10:22:33", "2006-01-02 15:04:05")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = ParseWithFormat("bad-input", "2006-01-02 15:04:05")
	if err == nil {
		t.Fatal("expected error for bad input, got none")
	}
}
