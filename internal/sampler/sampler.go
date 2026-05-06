package sampler

// Sampler provides line-based sampling of log output,
// allowing every Nth matching line to be retained.
type Sampler struct {
	step    int
	counter int
}

// New creates a Sampler that retains every nth line.
// A step of 0 or 1 means no sampling (all lines retained).
func New(step int) *Sampler {
	if step < 1 {
		step = 1
	}
	return &Sampler{step: step}
}

// Keep returns true if the current line should be retained
// based on the sampling step. It increments the internal
// counter on every call.
func (s *Sampler) Keep() bool {
	s.counter++
	if s.counter >= s.step {
		s.counter = 0
		return true
	}
	return false
}

// Reset resets the internal counter to zero.
func (s *Sampler) Reset() {
	s.counter = 0
}

// Step returns the configured sampling step.
func (s *Sampler) Step() int {
	return s.step
}

// Apply filters a slice of lines using the sampler,
// returning only the lines that pass the sampling criteria.
func (s *Sampler) Apply(lines []string) []string {
	s.Reset()
	out := make([]string, 0, len(lines)/s.step+1)
	for _, line := range lines {
		if s.Keep() {
			out = append(out, line)
		}
	}
	return out
}
