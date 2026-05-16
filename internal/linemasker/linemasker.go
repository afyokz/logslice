package linemasker

import "regexp"

// Masker replaces sensitive patterns in log lines with a fixed replacement string.
type Masker struct {
	enabled     bool
	patterns    []*regexp.Regexp
	replacement string
}

// New creates a Masker that replaces matches of any pattern in patterns with
// replacement. If patterns is empty, the Masker is disabled and Mask returns
// the original line unchanged. Returns an error if any pattern fails to compile.
func New(patterns []string, replacement string) (*Masker, error) {
	if len(patterns) == 0 {
		return &Masker{}, nil
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}
	if replacement == "" {
		replacement = "***"
	}
	return &Masker{
		enabled:     true,
		patterns:    compiled,
		replacement: replacement,
	}, nil
}

// Enabled reports whether the Masker will perform any substitutions.
func (m *Masker) Enabled() bool { return m.enabled }

// Replacement returns the string used to replace matched content.
func (m *Masker) Replacement() string { return m.replacement }

// Mask returns line with all pattern matches replaced by the replacement string.
// If the Masker is disabled, the original line is returned unchanged.
func (m *Masker) Mask(line string) string {
	if !m.enabled {
		return line
	}
	for _, re := range m.patterns {
		line = re.ReplaceAllString(line, m.replacement)
	}
	return line
}
