// Package linewindow provides a fixed-size sliding buffer that retains the N
// most recently observed log lines.
//
// It is intended to support "context lines" output — similar to grep's -B
// (before) flag — where a configurable number of lines preceding a match are
// emitted alongside the matching entry.
//
// Usage:
//
//	w := linewindow.New(3) // keep last 3 lines
//	for _, line := range allLines {
//		if matches(line) {
//			for _, ctx := range w.Lines() {
//				fmt.Println(ctx)
//			}
//			fmt.Println(line)
//		}
//		w.Push(line)
//	}
package linewindow
