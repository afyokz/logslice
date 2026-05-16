// Package linechunker provides a stateful line batcher that groups
// sequential log lines into fixed-size chunks.
//
// # Overview
//
// Chunker accumulates lines fed via [Chunker.Feed] and emits a complete
// batch ([]string) once the configured chunk size is reached. Partial
// batches left in the buffer can be retrieved at any time with
// [Chunker.Flush].
//
// When the chunk size is zero or negative the Chunker is disabled and
// every call to Feed returns immediately with a single-element slice,
// making it safe to use in pipelines regardless of configuration.
//
// # Typical usage
//
//	chunker := linechunker.New(500)
//	for _, line := range lines {
//		if batch := chunker.Feed(line); batch != nil {
//			process(batch)
//		}
//	}
//	if remainder := chunker.Flush(); remainder != nil {
//		process(remainder)
//	}
package linechunker
