package linerewriter

import "regexp"

// Rewriter replaces substrings in log lines using a list of regexp substitution rules.
type Rewriter struct {
	enabled bool
	rules   []rule
}

type rule struct {
	re          *regexp.Regexp
	replacement string
}

// Rule describes a single find-and-replace operation.
type Rule struct {
	Pattern     string
	Replacement string
}

// New creates a Rewriter from the given rules. Returns an error if any pattern
// fails to compile. When rules is empty the rewriter is disabled.
func New(rules []Rule) (*Rewriter, error) {
	if len(rules) == 0 {
		return &Rewriter{}, nil
	}
	compiled := make([]rule, 0, len(rules))
	for _, r := range rules {
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, rule{re: re, replacement: r.Replacement})
	}
	return &Rewriter{enabled: true, rules: compiled}, nil
}

// Enabled reports whether the rewriter has any active rules.
func (rw *Rewriter) Enabled() bool { return rw.enabled }

// Rewrite applies all substitution rules to line in order and returns the
// result. If the rewriter is disabled the original line is returned unchanged.
func (rw *Rewriter) Rewrite(line string) string {
	if !rw.enabled {
		return line
	}
	for _, r := range rw.rules {
		line = r.re.ReplaceAllString(line, r.replacement)
	}
	return line
}
