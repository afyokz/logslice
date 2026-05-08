// Package multilinecollector groups multi-line log entries into single
// logical records before they are passed to the rest of the logslice
// pipeline.
//
// Many log formats emit a primary header line followed by one or more
// continuation lines (e.g. Java stack traces, Python tracebacks, or
// wrapped SQL queries). Processing each physical line independently
// would break time-range filtering and pattern matching for those
// entries.
//
// A Collector wraps a [headerdetector.Detector] to decide which lines
// open a new record and which are continuations. Callers feed lines one
// at a time via Feed; when a new header is detected the previously
// accumulated record is returned. After the input is exhausted, Flush
// must be called to retrieve the final record.
package multilinecollector
