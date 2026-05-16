package lineclassifier

import "regexp"

// Classifier assigns a label to a log line based on named regexp patterns.
// If disabled (no patterns), Classify returns an empty string.
type Classifier struct {
	enabled  bool
	patterns []*entry
}

type entry struct {
	label string
	re    *regexp.Regexp
}

// New creates a Classifier from a map of label → pattern strings.
// Returns an error if any pattern fails to compile.
// If patterns is nil or empty the classifier is disabled.
func New(patterns map[string]string) (*Classifier, error) {
	if len(patterns) == 0 {
		return &Classifier{}, nil
	}
	entries := make([]*entry, 0, len(patterns))
	for label, pat := range patterns {
		re, err := regexp.Compile(pat)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &entry{label: label, re: re})
	}
	return &Classifier{enabled: true, patterns: entries}, nil
}

// Enabled reports whether the classifier has at least one pattern.
func (c *Classifier) Enabled() bool { return c.enabled }

// Classify returns the label of the first pattern that matches line.
// Returns an empty string when disabled or when no pattern matches.
func (c *Classifier) Classify(line string) string {
	if !c.enabled {
		return ""
	}
	for _, e := range c.patterns {
		if e.re.MatchString(line) {
			return e.label
		}
	}
	return ""
}

// Labels returns all configured label names in unspecified order.
func (c *Classifier) Labels() []string {
	out := make([]string, len(c.patterns))
	for i, e := range c.patterns {
		out[i] = e.label
	}
	return out
}
