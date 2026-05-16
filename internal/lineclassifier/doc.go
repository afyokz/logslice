// Package lineclassifier provides label-based classification of log lines.
//
// A Classifier holds a set of named regular-expression patterns. For each
// line it returns the label of the first matching pattern, or an empty
// string when no pattern matches.
//
// Example usage:
//
//	c, err := lineclassifier.New(map[string]string{
//		"error": `(?i)\berror\b`,
//		"warn":  `(?i)\bwarn\b`,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	label := c.Classify(line) // "error", "warn", or ""
//
// When constructed with a nil or empty pattern map the classifier is
// disabled and Classify always returns an empty string with no overhead.
package lineclassifier
