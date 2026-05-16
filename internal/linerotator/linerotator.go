package linerotator

import "io"

// Rotator writes lines to an underlying writer and rotates to a new writer
// after a configured number of lines have been written. When disabled (maxLines
// <= 0) all lines are forwarded to the single writer provided at construction.
type Rotator struct {
	current  io.Writer
	factory  func(seq int) (io.Writer, error)
	maxLines int
	written  int
	seq      int
	enabled  bool
}

// New creates a Rotator. factory is called each time a new segment is needed;
// seq starts at 0 and increments on every rotation. When maxLines <= 0 the
// rotator is disabled and factory is never called — all output goes to
// fallback.
func New(fallback io.Writer, factory func(seq int) (io.Writer, error), maxLines int) (*Rotator, error) {
	if maxLines <= 0 {
		return &Rotator{current: fallback, enabled: false}, nil
	}
	initial, err := factory(0)
	if err != nil {
		return nil, err
	}
	return &Rotator{
		current:  initial,
		factory:  factory,
		maxLines: maxLines,
		seq:      0,
		enabled:  true,
	}, nil
}

// WriteLine writes a single line (without a trailing newline) followed by '\n'
// to the current writer, rotating to a new writer when the line limit is hit.
func (r *Rotator) WriteLine(line string) error {
	if r.enabled && r.written > 0 && r.written%r.maxLines == 0 {
		if err := r.rotate(); err != nil {
			return err
		}
	}
	_, err := io.WriteString(r.current, line+"\n")
	if err != nil {
		return err
	}
	r.written++
	return nil
}

// Seq returns the index of the currently active segment (0-based).
func (r *Rotator) Seq() int { return r.seq }

// Written returns the total number of lines written across all segments.
func (r *Rotator) Written() int { return r.written }

func (r *Rotator) rotate() error {
	r.seq++
	w, err := r.factory(r.seq)
	if err != nil {
		return err
	}
	r.current = w
	return nil
}
