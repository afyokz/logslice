package lineformatter

import "fmt"

// Formatter applies optional transformations to a log line before output.
// It can prepend a line number and/or a custom prefix string.
type Formatter struct {
	numbered bool
	prefix   string
}

// New returns a Formatter. When numbered is true, Format prepends the
// 1-based sequence number of the line. prefix is prepended after the
// optional line number (ignored when empty).
func New(numbered bool, prefix string) *Formatter {
	return &Formatter{numbered: numbered, prefix: prefix}
}

// Format returns the formatted representation of line at position n
// (1-based). If neither numbering nor a prefix is configured the
// original line is returned unchanged.
func (f *Formatter) Format(n int, line string) string {
	switch {
	case f.numbered && f.prefix != "":
		return fmt.Sprintf("%d %s%s", n, f.prefix, line)
	case f.numbered:
		return fmt.Sprintf("%d %s", n, line)
	case f.prefix != "":
		return f.prefix + line
	default:
		return line
	}
}

// Numbered reports whether line numbering is enabled.
func (f *Formatter) Numbered() bool { return f.numbered }

// Prefix returns the configured prefix string.
func (f *Formatter) Prefix() string { return f.prefix }
