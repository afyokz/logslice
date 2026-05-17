package linerewriter_test

import (
	"fmt"

	"github.com/yourorg/logslice/internal/linerewriter"
)

func ExampleRewriter_Rewrite() {
	rw, err := linerewriter.New([]linerewriter.Rule{
		{Pattern: `\d+\.\d+\.\d+\.\d+`, Replacement: "<ip>"},
		{Pattern: `user=\S+`, Replacement: "user=<redacted>"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rw.Rewrite("connect from 192.168.1.10 user=alice"))
	// Output: connect from <ip> user=<redacted>
}

func ExampleRewriter_Rewrite_disabled() {
	rw, err := linerewriter.New(nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(rw.Rewrite("unchanged line"))
	// Output: unchanged line
}
