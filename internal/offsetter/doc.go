// Package offsetter provides timestamp shifting utilities for logslice.
//
// When processing log files from systems in different timezones, or when
// log timestamps need to be normalized to a common reference, Offsetter
// applies a fixed duration delta to time.Time values.
//
// Usage:
//
//	o := offsetter.New(-5 * time.Hour) // shift back 5 hours
//	shifted := o.Shift(parsedTimestamp)
//
// ShiftFrom and ShiftTo are convenience methods that treat zero-value
// time.Time as unbounded (i.e., no shift is applied), matching the
// semantics used by the scanner and pipeline packages.
package offsetter
