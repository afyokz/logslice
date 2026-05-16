package linemasker_test

import (
	"fmt"
	"log"

	"github.com/yourorg/logslice/internal/linemasker"
)

func ExampleMasker_Mask() {
	m, err := linemasker.New([]string{`password=\S+`}, "password=[REDACTED]")
	if err != nil {
		log.Fatal(err)
	}

	lines := []string{
		"2024-01-15T10:00:00Z user=alice password=hunter2 action=login",
		"2024-01-15T10:00:01Z user=alice action=view",
		"2024-01-15T10:00:02Z user=bob password=letmein action=login",
	}

	for _, l := range lines {
		fmt.Println(m.Mask(l))
	}

	// Output:
	// 2024-01-15T10:00:00Z user=alice password=[REDACTED] action=login
	// 2024-01-15T10:00:01Z user=alice action=view
	// 2024-01-15T10:00:02Z user=bob password=[REDACTED] action=login
}

func ExampleMasker_Mask_disabled() {
	// A Masker created with no patterns is a no-op.
	m, err := linemasker.New(nil, "")
	if err != nil {
		log.Fatal(err)
	}

	line := "password=topsecret"
	fmt.Println(m.Mask(line))

	// Output:
	// password=topsecret
}
