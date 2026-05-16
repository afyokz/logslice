package linesplitter

import "strings"

// Splitter splits a log line into named fields using a configurable delimiter
// and an ordered list of field names. It is disabled when no field names are
// provided, in which case Extract returns nil without allocating.
type Splitter struct {
	delimiter string
	fields    []string
	enabled   bool
	maxSplit  int
}

// New creates a new Splitter. delimiter defaults to a single space when empty.
// If fields is empty the splitter is disabled.
func New(delimiter string, fields []string) *Splitter {
	if delimiter == "" {
		delimiter = " "
	}
	enabled := len(fields) > 0
	return &Splitter{
		delimiter: delimiter,
		fields:    fields,
		enabled:   enabled,
		maxSplit:  len(fields),
	}
}

// Enabled reports whether the splitter is active.
func (s *Splitter) Enabled() bool { return s.enabled }

// Fields returns the configured field names.
func (s *Splitter) Fields() []string { return s.fields }

// Extract splits line by the configured delimiter and maps each part to the
// corresponding field name. If the line produces fewer parts than there are
// field names the missing fields are set to an empty string. Returns nil when
// the splitter is disabled.
func (s *Splitter) Extract(line string) map[string]string {
	if !s.enabled {
		return nil
	}
	parts := strings.SplitN(line, s.delimiter, s.maxSplit)
	out := make(map[string]string, len(s.fields))
	for i, name := range s.fields {
		if i < len(parts) {
			out[name] = parts[i]
		} else {
			out[name] = ""
		}
	}
	return out
}
