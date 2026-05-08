package tailreader

import (
	"io"
	"os"
	"time"
)

// TailReader reads new lines appended to a file, similar to `tail -f`.
type TailReader struct {
	file     *os.File
	pollInterval time.Duration
	stopCh   chan struct{}
}

// New opens the named file and returns a TailReader positioned at the end.
// pollInterval controls how often the file is polled for new data.
// If pollInterval is zero it defaults to 250ms.
func New(path string, pollInterval time.Duration) (*TailReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		f.Close()
		return nil, err
	}
	if pollInterval <= 0 {
		pollInterval = 250 * time.Millisecond
	}
	return &TailReader{
		file:         f,
		pollInterval: pollInterval,
		stopCh:       make(chan struct{}),
	}, nil
}

// Lines returns a channel that receives new lines as they are appended.
// The channel is closed when Stop is called or a read error occurs.
func (t *TailReader) Lines() <-chan string {
	ch := make(chan string, 64)
	go func() {
		defer close(ch)
		buf := make([]byte, 0, 4096)
		tmp := make([]byte, 4096)
		for {
			select {
			case <-t.stopCh:
				return
			default:
			}
			n, err := t.file.Read(tmp)
			if n > 0 {
				buf = append(buf, tmp[:n]...)
				for {
					idx := indexByte(buf, '\n')
					if idx < 0 {
						break
					}
					line := string(buf[:idx])
					buf = buf[idx+1:]
					ch <- line
				}
			}
			if err == io.EOF {
				time.Sleep(t.pollInterval)
				continue
			}
			if err != nil {
				return
			}
		}
	}()
	return ch
}

// Stop signals the tail goroutine to exit and closes the underlying file.
func (t *TailReader) Stop() error {
	close(t.stopCh)
	return t.file.Close()
}

// SeekToStart repositions the reader to the beginning of the file so that
// all existing lines will be re-emitted on the next call to Lines.
func (t *TailReader) SeekToStart() error {
	_, err := t.file.Seek(0, io.SeekStart)
	return err
}

func indexByte(b []byte, c byte) int {
	for i, v := range b {
		if v == c {
			return i
		}
	}
	return -1
}
