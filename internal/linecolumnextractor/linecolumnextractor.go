package linecolumnextractor

import (
	"strings"
)

// Extractor splits each log line into named columns using a fixed delimiter
// and returns a map of column name to value. If disabled or the line does
// not have enough fields, Extract returns nil.
type Extractor struct {
	enabled bool
	delimiter string
	columns []string
}

// New creates an Extractor that splits lines by delimiter and maps the
// resulting fields to the provided column names in order. If columns is
// empty the extractor is disabled. Delimiter defaults to a single space
// when the empty string is supplied.
func New(delimiter string, columns []string) *Extractor {
	if len(columns) == 0 {
		return &Extractor{}
	}
	if delimiter == "" {
		delimiter = " "
	}
	return &Extractor{
		enabled:   true,
		delimiter: delimiter,
		columns:   columns,
	}
}

// Enabled reports whether the extractor is active.
func (e *Extractor) Enabled() bool { return e.enabled }

// Columns returns the configured column names.
func (e *Extractor) Columns() []string { return e.columns }

// Extract splits line by the configured delimiter and returns a map of
// column name → field value. If the extractor is disabled or the line
// yields fewer fields than there are columns, Extract returns nil.
func (e *Extractor) Extract(line string) map[string]string {
	if !e.enabled {
		return nil
	}
	parts := strings.SplitN(line, e.delimiter, len(e.columns))
	if len(parts) < len(e.columns) {
		return nil
	}
	out := make(map[string]string, len(e.columns))
	for i, col := range e.columns {
		out[col] = parts[i]
	}
	return out
}
