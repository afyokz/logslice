// Package lineprefixer provides a Prefixer that prepends a fixed string
// and separator to every log line passing through the pipeline.
//
// It is useful for tagging lines with a source identifier, severity label,
// or any static annotation that should appear at the start of each line.
//
// Usage:
//
//	p := lineprefixer.New("[myapp]", " ")
//	formatted := p.Format(line) // "[myapp] original line"
//
// When the prefix is empty the Prefixer is disabled and Format returns
// the original line without any allocation.
package lineprefixer
