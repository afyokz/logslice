package linepadder

import "strings"

// Padder pads or truncates lines to a fixed width, optionally aligning content.
type Padder struct {
	enabled bool
	width   int
	align   Alignment
	pad     rune
}

// Alignment controls how content is positioned within the padded field.
type Alignment int

const (
	AlignLeft  Alignment = iota // default
	AlignRight
	AlignCenter
)

// New creates a Padder. Width <= 0 disables padding (lines pass through unchanged).
func New(width int, align Alignment, pad rune) *Padder {
	if pad == 0 {
		pad = ' '
	}
	return &Padder{
		enabled: width > 0,
		width:   width,
		align:   align,
		pad:     pad,
	}
}

// Enabled reports whether padding is active.
func (p *Padder) Enabled() bool { return p.enabled }

// Width returns the configured column width.
func (p *Padder) Width() int { return p.width }

// Pad returns the line padded (or truncated) to the configured width.
// If disabled, the original line is returned unchanged.
func (p *Padder) Pad(line string) string {
	if !p.enabled {
		return line
	}
	runes := []rune(line)
	l := len(runes)
	if l >= p.width {
		return string(runes[:p.width])
	}
	gap := p.width - l
	padStr := strings.Repeat(string(p.pad), gap)
	switch p.align {
	case AlignRight:
		return padStr + line
	case AlignCenter:
		left := gap / 2
		right := gap - left
		return strings.Repeat(string(p.pad), left) + line + strings.Repeat(string(p.pad), right)
	default: // AlignLeft
		return line + padStr
	}
}
