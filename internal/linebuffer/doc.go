// Package linebuffer provides a fixed-capacity ring buffer for log lines.
//
// Buffer retains the most recent N lines in insertion order, evicting the
// oldest entry when the capacity is exceeded. It is intended to capture
// trailing context around matched log entries during a pipeline scan.
//
// A zero or negative capacity disables the buffer; all operations become
// no-ops and Lines always returns nil, making it safe to use without
// conditional guards in calling code.
//
// Example:
//
//	buf := linebuffer.New(5)
//	for _, line := range logLines {
//		buf.Push(line)
//		if isMatch(line) {
//			context := buf.Lines() // last 5 lines including this one
//			process(context)
//			buf.Reset()
//		}
//	}
package linebuffer
