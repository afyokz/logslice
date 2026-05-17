// Package linescaler provides a numeric field scaler for log lines.
//
// A Scaler multiplies a single delimited field (identified by its 0-based
// column index) by a configurable factor.  This is useful when log metrics
// need to be converted between units — for example milliseconds to seconds,
// or bytes to kilobytes — during a pipeline export.
//
// # Disabling
//
// A factor of 0 or 1 disables the scaler; Scale returns the original line
// unchanged without any parsing overhead.
//
// # Errors
//
// Scale returns the original line together with a descriptive error when the
// target field index is out of range or the field value cannot be parsed as a
// floating-point number.  Callers may choose to skip, log, or propagate such
// errors depending on their pipeline requirements.
package linescaler
