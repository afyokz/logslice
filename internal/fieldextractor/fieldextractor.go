package fieldextractor

import (
	"regexp"
	"strings"
)

// Extractor extracts named fields from log lines using a configurable pattern.
type Extractor struct {
	pattern *regexp.Regexp
	fields  []string
	enabled bool
}

// New creates a new Extractor. If pattern is empty, the extractor is disabled.
// The pattern should use named capture groups, e.g. `(?P<level>\w+)`.
func New(pattern string) (*Extractor, error) {
	if pattern == "" {
		return &Extractor{enabled: false}, nil
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	fields := make([]string, 0)
	for _, name := range re.SubexpNames() {
		if name != "" {
			fields = append(fields, name)
		}
	}

	return &Extractor{
		pattern: re,
		fields:  fields,
		enabled: true,
	}, nil
}

// Enabled reports whether the extractor is active.
func (e *Extractor) Enabled() bool {
	return e.enabled
}

// Fields returns the list of named capture groups defined in the pattern.
func (e *Extractor) Fields() []string {
	return e.fields
}

// Extract parses the given line and returns a map of field name to value.
// Returns nil if the extractor is disabled or the pattern does not match.
func (e *Extractor) Extract(line string) map[string]string {
	if !e.enabled {
		return nil
	}

	match := e.pattern.FindStringSubmatch(line)
	if match == nil {
		return nil
	}

	result := make(map[string]string, len(e.fields))
	for i, name := range e.pattern.SubexpNames() {
		if name != "" && i < len(match) {
			result[name] = strings.TrimSpace(match[i])
		}
	}
	return result
}
