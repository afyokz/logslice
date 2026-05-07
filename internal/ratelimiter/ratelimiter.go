// Package ratelimiter provides a token-bucket style output rate limiter
// that restricts the number of log lines emitted per second.
package ratelimiter

import (
	"time"
)

// RateLimiter controls how many lines are allowed per second.
// A zero or negative rate means unlimited.
type RateLimiter struct {
	enabled  bool
	rate     int           // lines per second
	bucket   int           // current available tokens
	lastFill time.Time
	now      func() time.Time // injectable for testing
}

// New creates a RateLimiter that allows at most linesPerSec lines per second.
// If linesPerSec <= 0 the limiter is disabled and all lines are allowed.
func New(linesPerSec int) *RateLimiter {
	if linesPerSec <= 0 {
		return &RateLimiter{enabled: false, now: time.Now}
	}
	return &RateLimiter{
		enabled:  true,
		rate:     linesPerSec,
		bucket:   linesPerSec,
		lastFill: time.Now(),
		now:      time.Now,
	}
}

// Allow reports whether the next line should be emitted.
// It refills the token bucket based on elapsed time since the last call.
func (r *RateLimiter) Allow() bool {
	if !r.enabled {
		return true
	}

	now := r.now()
	elapsed := now.Sub(r.lastFill)

	if elapsed >= time.Second {
		// Refill tokens proportionally to elapsed seconds (cap at rate).
		secs := int(elapsed.Seconds())
		add := secs * r.rate
		r.bucket += add
		if r.bucket > r.rate {
			r.bucket = r.rate
		}
		r.lastFill = now
	}

	if r.bucket <= 0 {
		return false
	}
	r.bucket--
	return true
}

// Enabled reports whether rate limiting is active.
func (r *RateLimiter) Enabled() bool {
	return r.enabled
}

// Rate returns the configured lines-per-second limit.
func (r *RateLimiter) Rate() int {
	return r.rate
}
