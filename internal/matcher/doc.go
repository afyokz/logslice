// Package matcher provides pattern-based filtering of log lines using
// regular expressions.
//
// A Matcher can be configured with one or more include patterns and one or
// more exclude patterns:
//
//   - A line is excluded if it matches ANY exclude pattern.
//   - A line is included if it matches AT LEAST ONE include pattern (or if no
//     include patterns are configured, all lines pass).
//   - Exclude patterns always take priority over include patterns.
//
// Example usage:
//
//	cfg := matcher.Config{
//		IncludePatterns: []string{"ERROR", "FATAL"},
//		ExcludePatterns: []string{"healthcheck"},
//	}
//	m, err := matcher.New(cfg)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if m.Match(line) {
//		// process line
//	}
package matcher
