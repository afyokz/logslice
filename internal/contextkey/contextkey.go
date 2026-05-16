// Package contextkey defines typed keys for values stored in context.Context
// throughout the logslice pipeline. Using distinct types prevents key collisions
// between packages.
package contextkey

import (
	"context"
	"fmt"
)

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

// GetString retrieves a string value from ctx using the given key.
// It returns the value and true if found and the value is a string,
// or an empty string and false otherwise.
func GetString(ctx context.Context, k key) (string, bool) {
	v := ctx.Value(k)
	if v == nil {
		return "", false
	}
	s, ok := v.(string)
	if !ok {
		return "", false
	}
	return s, true
}

// MustGetString retrieves a string value from ctx using the given key.
// It panics if the key is not present or the value is not a string.
func MustGetString(ctx context.Context, k key) string {
	s, ok := GetString(ctx, k)
	if !ok {
		panic(fmt.Sprintf("contextkey: missing or invalid value for key %s", k))
	}
	return s
}
