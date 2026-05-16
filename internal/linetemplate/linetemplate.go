// Package linetemplate applies a Go text/template to structured log fields,
// allowing users to reformat matched log lines into a custom output string.
package linetemplate

import (
	"bytes"
	"text/template"
)

// Formatter renders log lines using a Go text/template applied to a map of
// named fields extracted from each line.
type Formatter struct {
	enabled bool
	tmpl    *template.Template
}

// New creates a Formatter for the given template string. If tmplStr is empty
// the formatter is disabled and Format returns the original line unchanged.
// Returns an error if the template fails to parse.
func New(tmplStr string) (*Formatter, error) {
	if tmplStr == "" {
		return &Formatter{enabled: false}, nil
	}
	t, err := template.New("line").Option("missingkey=zero").Parse(tmplStr)
	if err != nil {
		return nil, err
	}
	return &Formatter{enabled: true, tmpl: t}, nil
}

// Enabled reports whether the formatter will rewrite lines.
func (f *Formatter) Enabled() bool { return f.enabled }

// Format applies the template to fields and returns the rendered string.
// If the formatter is disabled, the original line is returned unchanged.
// If template execution fails the original line is returned unchanged.
func (f *Formatter) Format(line string, fields map[string]string) string {
	if !f.enabled {
		return line
	}
	data := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		data[k] = v
	}
	var buf bytes.Buffer
	if err := f.tmpl.Execute(&buf, data); err != nil {
		return line
	}
	return buf.String()
}
