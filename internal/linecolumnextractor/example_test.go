package linecolumnextractor_test

import (
	"fmt"

	"github.com/yourorg/logslice/internal/linecolumnextractor"
)

func ExampleExtractor_Extract() {
	e := linecolumnextractor.New(" ", []string{"timestamp", "level", "message"})

	fields := e.Extract("2024-06-01T12:00:00Z INFO server started successfully")
	if fields == nil {
		fmt.Println("no match")
		return
	}
	fmt.Println(fields["timestamp"])
	fmt.Println(fields["level"])
	fmt.Println(fields["message"])
	// Output:
	// 2024-06-01T12:00:00Z
	// INFO
	// server started successfully
}

func ExampleExtractor_Extract_csv() {
	e := linecolumnextractor.New(",", []string{"ip", "method", "path", "status"})

	fields := e.Extract("192.168.1.1,GET,/health,200")
	if fields == nil {
		fmt.Println("no match")
		return
	}
	fmt.Printf("ip=%s method=%s path=%s status=%s\n",
		fields["ip"], fields["method"], fields["path"], fields["status"])
	// Output:
	// ip=192.168.1.1 method=GET path=/health status=200
}
