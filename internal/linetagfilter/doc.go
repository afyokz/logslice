// Package linetagfilter filters log lines by structured key=value tags.
//
// A Filter is constructed with a list of tag expressions:
//
//	- "key=value"   — exact match
//	- "key=~pattern" — regular-expression match against the value
//
// All expressions must be satisfied for a line to be kept (AND semantics).
// When no expressions are provided the filter is disabled and every line
// passes through unchanged.
//
// Example:
//
//	f, err := linetagfilter.New([]string{"level=error", "svc=~^(auth|api)$"})
//	if err != nil { /* handle */ }
//	if f.Keep(line) {
//	    // emit line
//	}
package linetagfilter
