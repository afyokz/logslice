package linetagfilter_test

import (
	"fmt"

	"github.com/user/logslice/internal/linetagfilter"
)

func ExampleFilter_Keep_exact() {
	f, _ := linetagfilter.New([]string{"level=error"})

	lines := []string{
		`ts=1 level=info  msg="starting up"`,
		`ts=2 level=error msg="disk full" svc=storage`,
		`ts=3 level=warn  msg="high mem"`,
	}

	for _, l := range lines {
		if f.Keep(l) {
			fmt.Println(l)
		}
	}
	// Output:
	// ts=2 level=error msg="disk full" svc=storage
}

func ExampleFilter_Keep_regex() {
	f, _ := linetagfilter.New([]string{"level=~^(error|warn)$", "svc=auth"})

	lines := []string{
		`ts=1 level=info  svc=auth  msg="login ok"`,
		`ts=2 level=error svc=auth  msg="bad token"`,
		`ts=3 level=warn  svc=api   msg="slow"`,
		`ts=4 level=warn  svc=auth  msg="rate limit"`,
	}

	for _, l := range lines {
		if f.Keep(l) {
			fmt.Println(l)
		}
	}
	// Output:
	// ts=2 level=error svc=auth  msg="bad token"
	// ts=4 level=warn  svc=auth  msg="rate limit"
}
