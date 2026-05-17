// Package linetagfilter provides filtering of log lines based on structured
// key=value tags embedded in each line.
package linetagfilter

import (
	"regexp"
	"strings"
)

// Filter retains or rejects lines based on tag key=value pairs.
type Filter struct {
	enabled  bool
	required map[string]*regexp.Regexp // key → value pattern (nil means any value)
}

// New creates a Filter from a slice of tag expressions in the form "key=value"
// or "key=~pattern" (regex match). An empty tags slice disables filtering.
// Returns an error if any regex pattern is invalid.
func New(tags []string) (*Filter, error) {
	if len(tags) == 0 {
		return &Filter{}, nil
	}
	required := make(map[string]*regexp.Regexp, len(tags))
	for _, t := range tags {
		key, val, found := strings.Cut(t, "=")
		if !found || key == "" {
			continue
		}
		if strings.HasPrefix(val, "~") {
			re, err := regexp.Compile(val[1:])
			if err != nil {
				return nil, err
			}
			required[key] = re
		} else {
			// exact match stored as anchored literal regex
			re := regexp.MustCompile("^" + regexp.QuoteMeta(val) + "$")
			required[key] = re
		}
	}
	return &Filter{enabled: len(required) > 0, required: required}, nil
}

// Enabled reports whether the filter is active.
func (f *Filter) Enabled() bool { return f.enabled }

// Keep returns true when the line satisfies all required tag constraints.
// If the filter is disabled, Keep always returns true.
func (f *Filter) Keep(line string) bool {
	if !f.enabled {
		return true
	}
	matched := make(map[string]bool, len(f.required))
	for _, field := range strings.Fields(line) {
		k, v, ok := strings.Cut(field, "=")
		if !ok {
			continue
		}
		if re, want := f.required[k]; want {
			if re.MatchString(v) {
				matched[k] = true
			}
		}
	}
	return len(matched) == len(f.required)
}
