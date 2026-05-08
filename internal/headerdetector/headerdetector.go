package headerdetector

import (
	"regexp"
	"strings"
)

// Detector identifies whether a log line begins with a recognisable
// timestamp header so that continuation lines (stack traces, wrapped
// text, etc.) can be handled correctly by the pipeline.
type Detector struct {
	enabled  bool
	patterns []*regexp.Regexp
}

// defaultPatterns covers the most common timestamp prefixes found in
// real-world log files.
var defaultPatterns = []string{
	`^\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}`,   // ISO-8601 / RFC-3339
	`^\d{2}/\d{2}/\d{4} \d{2}:\d{2}:\d{2}`,        // US date style
	`^\[\d{4}-\d{2}-\d{2}`,                         // bracketed date
	`^\w{3} [ \d]\d \d{2}:\d{2}:\d{2}`,             // syslog style
}

// New returns a Detector compiled from the supplied pattern strings.
// If patterns is empty the built-in defaults are used. An error is
// returned when any pattern fails to compile.
func New(patterns []string) (*Detector, error) {
	if len(patterns) == 0 {
		patterns = defaultPatterns
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}
	return &Detector{enabled: true, patterns: compiled}, nil
}

// IsHeader reports whether line starts a new log record, i.e. it
// matches at least one of the compiled timestamp patterns.
func (d *Detector) IsHeader(line string) bool {
	if !d.enabled {
		return true // treat every line as a header when disabled
	}
	for _, re := range d.patterns {
		if re.MatchString(line) {
			return true
		}
	}
	return false
}

// IsContinuation is the inverse of IsHeader.
func (d *Detector) IsContinuation(line string) bool {
	return !d.IsHeader(line)
}

// StripTimestamp removes the leading timestamp portion from line and
// returns the remainder trimmed of surrounding whitespace. If no
// pattern matches, the original line is returned unchanged.
func (d *Detector) StripTimestamp(line string) string {
	for _, re := range d.patterns {
		loc := re.FindStringIndex(line)
		if loc != nil {
			return strings.TrimSpace(line[loc[1]:])
		}
	}
	return line
}
