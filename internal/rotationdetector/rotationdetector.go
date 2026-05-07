package rotationdetector

import (
	"os"
	"sync"
)

// Detector watches a log file for rotation events by tracking its inode
// and size. It reports when the underlying file has been replaced or truncated.
type Detector struct {
	mu      sync.Mutex
	path    string
	inode   uint64
	size    int64
	enabled bool
}

// New creates a Detector for the given file path. If path is empty or
// the file cannot be stat'd at construction time, detection is disabled.
func New(path string) (*Detector, error) {
	if path == "" {
		return &Detector{enabled: false}, nil
	}
	d := &Detector{path: path, enabled: true}
	if err := d.snapshot(); err != nil {
		return nil, err
	}
	return d, nil
}

// Rotated returns true if the file has been rotated (replaced or truncated)
// since the last call to Reset, or since construction. It is safe for
// concurrent use.
func (d *Detector) Rotated() (bool, error) {
	if !d.enabled {
		return false, nil
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	fi, err := os.Stat(d.path)
	if err != nil {
		// File disappeared — treat as rotated.
		return true, nil
	}
	currentInode := inode(fi)
	currentSize := fi.Size()

	if currentInode != d.inode || currentSize < d.size {
		return true, nil
	}
	return false, nil
}

// Reset updates the baseline snapshot to the current file state.
func (d *Detector) Reset() error {
	if !d.enabled {
		return nil
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.snapshot()
}

// Enabled reports whether rotation detection is active.
func (d *Detector) Enabled() bool {
	return d.enabled
}

func (d *Detector) snapshot() error {
	fi, err := os.Stat(d.path)
	if err != nil {
		return err
	}
	d.inode = inode(fi)
	d.size = fi.Size()
	return nil
}
