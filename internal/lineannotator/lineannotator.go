// Package lineannotator prepends a configurable prefix tag to each log line,
// useful for labelling output from multiple sources or adding metadata markers.
package lineannotator

import "fmt"

// Annotator prepends a tag to log lines.
type Annotator struct {
	enabled bool
	tag     string
	sep     string
}

// New creates an Annotator that prepends tag+sep before each line.
// If tag is empty, the annotator is disabled and lines pass through unchanged.
// sep defaults to ": " when empty.
func New(tag, sep string) *Annotator {
	if tag == "" {
		return &Annotator{enabled: false}
	}
	if sep == "" {
		sep = ": "
	}
	return &Annotator{
		enabled: true,
		tag:     tag,
		sep:     sep,
	}
}

// Enabled reports whether annotation is active.
func (a *Annotator) Enabled() bool {
	return a.enabled
}

// Tag returns the configured annotation tag.
func (a *Annotator) Tag() string {
	return a.tag
}

// Annotate returns the line with the tag prepended.
// If the annotator is disabled, the original line is returned unchanged.
func (a *Annotator) Annotate(line string) string {
	if !a.enabled {
		return line
	}
	return fmt.Sprintf("%s%s%s", a.tag, a.sep, line)
}
