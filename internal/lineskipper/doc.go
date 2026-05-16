// Package lineskipper implements a line skipper that discards the first N
// lines of a log stream.
//
// This is useful when log files begin with static headers, version banners,
// or other preamble content that should be excluded from time-range slicing.
//
// Usage:
//
//	s := lineskipper.New(5) // skip the first 5 lines
//	for _, line := range lines {
//		if s.Keep(line) {
//			// process line
//		}
//	}
//
// A zero or negative argument disables the skipper, making Keep always
// return true with no overhead.
package lineskipper
