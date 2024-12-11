package backoff

import (
	"time"
)

const (
	defaultStep = 2
)

type Backoff struct {
	min        time.Duration
	max        time.Duration
	maxAttempt int
	attemptNum int
	nextDelay  time.Duration
}

func New(min, max time.Duration, maxAttempt int) *Backoff {
	return &Backoff{
		min:        min,
		max:        max,
		maxAttempt: maxAttempt,
		nextDelay:  min,
	}
}

const Stop time.Duration = -1

func (b *Backoff) Next() time.Duration {
	if b.attemptNum >= b.maxAttempt {
		return Stop
	}
	b.attemptNum++
	delay := min(b.nextDelay, b.max)
	b.nextDelay += defaultStep * time.Second
	return delay
}

func (b *Backoff) Reset() {
	b.attemptNum = 0
	b.nextDelay = b.min
}
