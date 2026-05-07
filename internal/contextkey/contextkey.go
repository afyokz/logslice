// Package contextkey defines typed keys for values stored in context.Context
// throughout the logslice pipeline. Using distinct types prevents key collisions
// between packages.
package contextkey

// key is an unexported type for context keys in this package.
type key int

const (
	// JobID identifies a unique processing job or request.
	JobID key = iota

	// InputFile is the path of the log file being processed.
	InputFile

	// Verbose controls whether verbose/debug output is enabled.
	Verbose
)

// String returns a human-readable name for the context key.
func (k key) String() string {
	switch k {
	case JobID:
		return "JobID"
	case InputFile:
		return "InputFile"
	case Verbose:
		return "Verbose"
	default:
		return "unknown"
	}
}
