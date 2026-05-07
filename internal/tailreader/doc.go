// Package tailreader provides a live-tail reader that monitors a log file
// for newly appended lines and delivers them over a channel.
//
// It is designed to integrate with the logslice pipeline for streaming
// (follow) mode, analogous to `tail -f`. The reader polls the underlying
// file at a configurable interval and emits complete newline-terminated
// lines as strings.
//
// Usage:
//
//	tr, err := tailreader.New("/var/log/app.log", 100*time.Millisecond)
//	if err != nil { ... }
//	defer tr.Stop()
//	for line := range tr.Lines() {
//	    fmt.Println(line)
//	}
package tailreader
