package linecapper

// Capper limits the byte length of a line, hard-capping it at a maximum
// number of bytes and appending an optional suffix (e.g. "...") when the
// line is truncated.
//
// When MaxBytes is zero or negative the Capper is disabled and Cap returns
// the original string unchanged.
type Capper struct {
	maxBytes int
	suffix   string
	enabled  bool
}

// New returns a Capper that caps lines at maxBytes bytes.
// suffix is appended to capped lines; pass an empty string for no suffix.
// If maxBytes <= 0 the Capper is disabled.
func New(maxBytes int, suffix string) *Capper {
	if maxBytes <= 0 {
		return &Capper{}
	}
	return &Capper{
		maxBytes: maxBytes,
		suffix:   suffix,
		enabled:  true,
	}
}

// Enabled reports whether the Capper is active.
func (c *Capper) Enabled() bool { return c.enabled }

// MaxBytes returns the configured byte cap (0 when disabled).
func (c *Capper) MaxBytes() int { return c.maxBytes }

// Cap returns line unchanged when the Capper is disabled or the line already
// fits within maxBytes. Otherwise it truncates the line so that the result
// (including the suffix) is at most maxBytes bytes long.
//
// If the suffix itself is longer than maxBytes, Cap returns only the suffix
// truncated to maxBytes.
func (c *Capper) Cap(line string) string {
	if !c.enabled {
		return line
	}
	if len(line) <= c.maxBytes {
		return line
	}
	suffixBytes := len(c.suffix)
	keep := c.maxBytes - suffixBytes
	if keep <= 0 {
		// suffix alone fills the budget
		if suffixBytes > c.maxBytes {
			return c.suffix[:c.maxBytes]
		}
		return c.suffix
	}
	return line[:keep] + c.suffix
}
