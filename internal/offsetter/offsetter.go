package offsetter

import "time"

// Offsetter shifts timestamps by a fixed duration, useful for normalizing
// log files from different timezones or adjusting relative timestamps.
type Offsetter struct {
	offset time.Duration
}

// New creates a new Offsetter with the given duration offset.
// A positive offset shifts timestamps forward; negative shifts backward.
func New(offset time.Duration) *Offsetter {
	return &Offsetter{offset: offset}
}

// Shift applies the configured offset to the given timestamp.
func (o *Offsetter) Shift(t time.Time) time.Time {
	return t.Add(o.offset)
}

// ShiftFrom applies the offset to the from timestamp if non-zero.
func (o *Offsetter) ShiftFrom(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	return t.Add(o.offset)
}

// ShiftTo applies the offset to the to timestamp if non-zero.
func (o *Offsetter) ShiftTo(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	return t.Add(o.offset)
}

// IsZero reports whether the offset is zero (no-op).
func (o *Offsetter) IsZero() bool {
	return o.offset == 0
}

// String returns a human-readable representation of the offset.
func (o *Offsetter) String() string {
	if o.offset == 0 {
		return "no offset"
	}
	return o.offset.String()
}
