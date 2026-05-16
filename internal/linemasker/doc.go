// Package linemasker provides a Masker that redacts sensitive content in log
// lines by replacing substrings that match one or more regular expressions with
// a configurable replacement string.
//
// Typical use cases include masking passwords, API tokens, credit-card numbers,
// or any other PII before the line is written to an output stream.
//
// Example:
//
//	m, err := linemasker.New([]string{`password=\S+`, `token=\S+`}, "***")
//	if err != nil {
//		log.Fatal(err)
//	}
//	safe := m.Mask(rawLine)
package linemasker
