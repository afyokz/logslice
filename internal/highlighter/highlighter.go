// Package highlighter provides optional ANSI color highlighting
// for matched terms within log lines.
package highlighter

import (
	"regexp"
	"strings"
)

const (
	ansiReset  = "\033[0m"
	ansiYellow = "\033[33m"
	ansiRed    = "\033[31m"
	ansiBold   = "\033[1m"
)

// Highlighter wraps matched substrings in ANSI escape codes.
type Highlighter struct {
	patterns []*regexp.Regexp
	color    string
	enabled  bool
}

// New creates a Highlighter for the given keyword terms.
// If no terms are provided or enabled is false, Highlight is a no-op.
func New(enabled bool, color string, terms []string) (*Highlighter, error) {
	h := &Highlighter{enabled: enabled, color: resolveColor(color)}
	if !enabled {
		return h, nil
	}
	for _, t := range terms {
		if t == "" {
			continue
		}
		re, err := regexp.Compile(`(?i)` + regexp.QuoteMeta(t))
		if err != nil {
			return nil, err
		}
		h.patterns = append(h.patterns, re)
	}
	return h, nil
}

// Highlight returns line with all matched terms wrapped in ANSI color codes.
// If the highlighter is disabled or has no patterns, line is returned unchanged.
func (h *Highlighter) Highlight(line string) string {
	if !h.enabled || len(h.patterns) == 0 {
		return line
	}
	for _, re := range h.patterns {
		line = re.ReplaceAllStringFunc(line, func(match string) string {
			return h.color + ansiBold + match + ansiReset
		})
	}
	return line
}

// Enabled reports whether highlighting is active.
func (h *Highlighter) Enabled() bool {
	return h.enabled
}

// PatternCount returns the number of compiled highlight patterns.
func (h *Highlighter) PatternCount() int {
	return len(h.patterns)
}

func resolveColor(color string) string {
	switch strings.ToLower(color) {
	case "red":
		return ansiRed
	case "yellow", "":
		return ansiYellow
	default:
		return ansiYellow
	}
}
