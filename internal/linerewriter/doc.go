// Package linerewriter provides a Rewriter that applies a sequence of
// regular-expression substitution rules to log lines.
//
// Rules are compiled once at construction time and applied in the order they
// were provided. When the rule list is empty the rewriter is disabled and
// every call to Rewrite returns the original line unchanged, with no
// allocation overhead.
//
// Example usage:
//
//	rw, err := linerewriter.New([]linerewriter.Rule{
//		{Pattern: `\d{4}-\d{2}-\d{2}`, Replacement: "DATE"},
//		{Pattern: `user=\S+`,          Replacement: "user=REDACTED"},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(rw.Rewrite("2024-01-15 user=alice logged in"))
//	// Output: DATE user=REDACTED logged in
package linerewriter
