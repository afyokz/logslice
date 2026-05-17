// Package linecensor provides a Censor type that redacts substrings of log
// lines matching one or more regular-expression patterns.
//
// Typical use-cases include removing passwords, API keys, PII (personally
// identifiable information), or any other sensitive data before the lines are
// exported or displayed.
//
// Usage:
//
//	c, err := linecensor.New([]string{`password=\S+`, `token=\S+`}, "***")
//	if err != nil {
//		log.Fatal(err)
//	}
//	cleaned := c.Apply(rawLine)
//
// When the patterns slice is empty the Censor is disabled and Apply returns
// the original line unchanged with zero allocations.
package linecensor
