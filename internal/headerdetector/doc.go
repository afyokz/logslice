// Package headerdetector identifies whether a log line begins a new
// log record (a "header" line) or is a continuation of the previous
// record (e.g. a stack-trace frame or a wrapped message body).
//
// Detection is performed by matching the start of each line against a
// set of compiled regular expressions. A library of common timestamp
// formats is provided as the default set; callers may supply their own
// patterns when constructing a Detector.
//
// Typical usage:
//
//	det, err := headerdetector.New(nil) // use built-in patterns
//	if err != nil { ... }
//
//	for _, line := range lines {
//		if det.IsHeader(line) {
//			// start of a new log record
//		} else {
//			// continuation – attach to previous record
//		}
//	}
package headerdetector
