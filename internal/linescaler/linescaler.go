package linescaler

import (
	"fmt"
	"strconv"
	"strings"
)

// Scaler normalises numeric fields found in a log line by applying a
// multiplicative scale factor.  When disabled (factor == 1 or == 0) every
// line is returned unchanged.
type Scaler struct {
	enabled bool
	factor  float64
	field   int // 0-based column index
	sep     string
}

// New creates a Scaler that multiplies the value at column fieldIndex
// (0-based) by factor using sep as the field delimiter.
// A factor of 0 or 1 disables scaling.
func New(factor float64, fieldIndex int, sep string) *Scaler {
	if sep == "" {
		sep = " "
	}
	enabled := factor != 0 && factor != 1
	return &Scaler{
		enabled: enabled,
		factor:  factor,
		field:   fieldIndex,
		sep:     sep,
	}
}

// Enabled reports whether scaling is active.
func (s *Scaler) Enabled() bool { return s.enabled }

// Factor returns the configured scale factor.
func (s *Scaler) Factor() float64 { return s.factor }

// Scale applies the scale factor to the target field of line.
// If the field cannot be parsed as a float the original line is returned
// together with an error; other fields are never modified.
func (s *Scaler) Scale(line string) (string, error) {
	if !s.enabled {
		return line, nil
	}
	parts := strings.Split(line, s.sep)
	if s.field >= len(parts) {
		return line, fmt.Errorf("linescaler: field index %d out of range (got %d fields)", s.field, len(parts))
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(parts[s.field]), 64)
	if err != nil {
		return line, fmt.Errorf("linescaler: cannot parse %q as float: %w", parts[s.field], err)
	}
	parts[s.field] = strconv.FormatFloat(v*s.factor, 'f', -1, 64)
	return strings.Join(parts, s.sep), nil
}
