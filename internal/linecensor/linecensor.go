package linecensor

import (
	"regexp"
	"strings"
)

// Censor replaces matched substrings with a fixed placeholder string.
// It is disabled when no patterns are provided.
type Censor struct {
	enabled     bool
	patterns    []*regexp.Regexp
	placeholder string
}

// New creates a Censor that replaces each match of any pattern in patterns
// with placeholder. Returns an error if any pattern fails to compile.
// If patterns is empty the Censor is disabled and Censor returns the line
// unchanged.
func New(patterns []string, placeholder string) (*Censor, error) {
	if len(patterns) == 0 {
		return &Censor{}, nil
	}
	if placeholder == "" {
		placeholder = "***"
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}
	return &Censor{
		enabled:     true,
		patterns:    compiled,
		placeholder: placeholder,
	}, nil
}

// Enabled reports whether the censor is active.
func (c *Censor) Enabled() bool { return c.enabled }

// Placeholder returns the replacement string used for matched substrings.
func (c *Censor) Placeholder() string { return c.placeholder }

// Apply censors line by replacing every match of every configured pattern
// with the placeholder. Returns line unchanged when disabled.
func (c *Censor) Apply(line string) string {
	if !c.enabled {
		return line
	}
	for _, re := range c.patterns {
		line = re.ReplaceAllString(line, c.placeholder)
	}
	return line
}

// ContainsAny reports whether line matches at least one of the configured
// patterns. Always returns false when disabled.
func (c *Censor) ContainsAny(line string) bool {
	if !c.enabled {
		return false
	}
	for _, re := range c.patterns {
		if re.MatchString(line) {
			return true
		}
	}
	return false
}

// ApplyAll censors every line in lines, returning a new slice.
func (c *Censor) ApplyAll(lines []string) []string {
	if !c.enabled {
		return lines
	}
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = c.Apply(l)
	}
	return out
}

// String returns a human-readable summary of the censor configuration.
func (c *Censor) String() string {
	if !c.enabled {
		return "Censor(disabled)"
	}
	return "Censor(patterns=" + strings.Join(patternStrings(c.patterns), ",") + ", placeholder=" + c.placeholder + ")"
}

func patternStrings(res []*regexp.Regexp) []string {
	s := make([]string, len(res))
	for i, re := range res {
		s[i] = re.String()
	}
	return s
}
