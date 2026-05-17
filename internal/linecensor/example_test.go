package linecensor_test

import (
	"fmt"

	"github.com/yourorg/logslice/internal/linecensor"
)

func ExampleCensor_Apply() {
	c, err := linecensor.New([]string{`password=\S+`}, "password=[REDACTED]")
	if err != nil {
		panic(err)
	}
	line := "2024-01-02 login user=alice password=hunter2 ok"
	fmt.Println(c.Apply(line))
	// Output: 2024-01-02 login user=alice password=[REDACTED] ok
}

func ExampleCensor_Apply_disabled() {
	c, err := linecensor.New(nil, "***")
	if err != nil {
		panic(err)
	}
	line := "nothing to censor here"
	fmt.Println(c.Apply(line))
	// Output: nothing to censor here
}
