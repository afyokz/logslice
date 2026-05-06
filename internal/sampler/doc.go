// Package sampler provides line-based sampling for log output.
//
// When processing large log files, it is often useful to reduce
// the volume of output by retaining only every Nth matching line.
// This package implements a simple counter-based sampler that can
// be applied to any stream of log lines after filtering.
//
// Example usage:
//
//	s := sampler.New(10) // keep every 10th line
//	for _, line := range lines {
//		if s.Keep() {
//			fmt.Println(line)
//		}
//	}
//
// The sampler is stateful; call Reset() to reuse it across
// multiple passes over the same data.
package sampler
