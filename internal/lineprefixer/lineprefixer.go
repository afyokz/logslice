package lineprefixer

import "strings"

// Prefixer prepends a fixed string to each line.
// When disabled (empty prefix), Format returns the original line unchanged.
type Prefixer struct {
	prefix    string
	separator string
	enabled   bool
}

// New creates a Prefixer that prepends prefix + separator to each line.
// If prefix is empty, the Prefixer is disabled and Format is a no-op.
// If separator is empty it defaults to a single space.
func New(prefix, separator string) *Prefixer {
	if prefix == "" {
		return &Prefixer{}
	}
	sep := separator
	if sep == "" {
		sep = " "
	}
	return &Prefixer{
		prefix:    prefix,
		separator: sep,
		enabled:   true,
	}
}

// Enabled reports whether the Prefixer will modify lines.
func (p *Prefixer) Enabled() bool { return p.enabled }

// Prefix returns the configured prefix string.
func (p *Prefixer) Prefix() string { return p.prefix }

// Separator returns the separator placed between the prefix and the line.
func (p *Prefixer) Separator() string { return p.separator }

// Format prepends the prefix and separator to line.
// If the Prefixer is disabled, line is returned unchanged.
func (p *Prefixer) Format(line string) string {
	if !p.enabled {
		return line
	}
	var b strings.Builder
	b.Grow(len(p.prefix) + len(p.separator) + len(line))
	b.WriteString(p.prefix)
	b.WriteString(p.separator)
	b.WriteString(line)
	return b.String()
}
