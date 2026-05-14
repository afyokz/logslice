package linemerger

import "strings"

// Merger joins adjacent log lines that share a common prefix (e.g. same
// timestamp prefix) into a single logical line separated by a configurable
// delimiter.
type Merger struct {
	enabled   bool
	delimiter string
	prefix    int // number of leading bytes used as the grouping key
	buf       []string
	currentKey string
}

// New creates a Merger. prefixLen is the number of leading characters used to
// identify lines that belong to the same group. A delimiter is inserted
// between joined lines. When prefixLen <= 0 the merger is disabled.
func New(prefixLen int, delimiter string) *Merger {
	if delimiter == "" {
		delimiter = " "
	}
	return &Merger{
		enabled:   prefixLen > 0,
		delimiter: delimiter,
		prefix:    prefixLen,
	}
}

// Enabled reports whether the merger is active.
func (m *Merger) Enabled() bool { return m.enabled }

// Feed adds a line to the merger. If the line belongs to the current group it
// is buffered and nil is returned. When a new group starts the previously
// buffered group is returned as a merged string.
func (m *Merger) Feed(line string) (string, bool) {
	if !m.enabled {
		return line, true
	}
	key := m.keyOf(line)
	if key != m.currentKey && len(m.buf) > 0 {
		merged := strings.Join(m.buf, m.delimiter)
		m.buf = m.buf[:0]
		m.currentKey = key
		m.buf = append(m.buf, line)
		return merged, true
	}
	m.currentKey = key
	m.buf = append(m.buf, line)
	return "", false
}

// Flush returns any remaining buffered lines as a single merged string.
// Returns empty string and false when there is nothing buffered.
func (m *Merger) Flush() (string, bool) {
	if len(m.buf) == 0 {
		return "", false
	}
	merged := strings.Join(m.buf, m.delimiter)
	m.buf = m.buf[:0]
	m.currentKey = ""
	return merged, true
}

// Reset discards all buffered state.
func (m *Merger) Reset() {
	m.buf = m.buf[:0]
	m.currentKey = ""
}

func (m *Merger) keyOf(line string) string {
	if len(line) <= m.prefix {
		return line
	}
	return line[:m.prefix]
}
