// Package limiter provides a simple line-count limiter for controlling
// the maximum number of output lines produced by logslice.
//
// A Limiter is created with a maximum line count. Each call to Allow
// increments an internal counter and returns true until the limit is
// reached, after which it returns false. A limit of zero or any negative
// value disables enforcement, allowing unlimited lines through.
//
// Typical usage:
//
//	l := limiter.New(cfg.MaxLines)
//	for _, line := range lines {
//		if !l.Allow() {
//			break
//		}
//		fmt.Fprintln(out, line)
//	}
package limiter
