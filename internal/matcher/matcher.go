// Package matcher provides pattern-based line filtering for log entries.
package matcher

import (
	"regexp"
	"strings"
)

// Matcher filters log lines based on include/exclude patterns.
type Matcher struct {
	include []*regexp.Regexp
	exclude []*regexp.Regexp
}

// Config holds the configuration for a Matcher.
type Config struct {
	// IncludePatterns are regex patterns; a line must match at least one to be included.
	IncludePatterns []string
	// ExcludePatterns are regex patterns; a line matching any will be excluded.
	ExcludePatterns []string
}

// New creates a new Matcher from the given Config.
// Returns an error if any pattern fails to compile.
func New(cfg Config) (*Matcher, error) {
	m := &Matcher{}

	for _, p := range cfg.IncludePatterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		m.include = append(m.include, re)
	}

	for _, p := range cfg.ExcludePatterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		m.exclude = append(m.exclude, re)
	}

	return m, nil
}

// Match returns true if the line passes include and exclude filters.
// If no include patterns are set, all lines pass the include check.
func (m *Matcher) Match(line string) bool {
	for _, re := range m.exclude {
		if re.MatchString(line) {
			return false
		}
	}

	if len(m.include) == 0 {
		return true
	}

	for _, re := range m.include {
		if re.MatchString(line) {
			return true
		}
	}

	return false
}

// MatchAll returns true only if the line matches ALL include patterns.
func (m *Matcher) MatchAll(line string) bool {
	for _, re := range m.exclude {
		if re.MatchString(line) {
			return false
		}
	}

	for _, re := range m.include {
		if !re.MatchString(line) {
			return false
		}
	}

	return true
}

// ContainsAny is a lightweight helper that checks if line contains any of the
// given literal substrings without compiling a regex.
func ContainsAny(line string, terms []string) bool {
	for _, t := range terms {
		if strings.Contains(line, t) {
			return true
		}
	}
	return false
}
