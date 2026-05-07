//go:build !windows

package rotationdetector

import (
	"io/fs"
	"syscall"
)

// inode extracts the inode number from a FileInfo on Unix-like systems.
func inode(fi fs.FileInfo) uint64 {
	if stat, ok := fi.Sys().(*syscall.Stat_t); ok {
		return stat.Ino
	}
	return 0
}
