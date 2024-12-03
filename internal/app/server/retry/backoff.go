package retry

import (
	"time"
)

const (
	DefaultDurFactor = 2 * time.Second
)

type FnBackoff func(attemptNum int, min, max time.Duration) time.Duration

type Backoff struct {
	backoff    FnBackoff
	min        time.Duration
	max        time.Duration
	maxAttempt int
	attemptNum int
}

func NewBackoff(min, max time.Duration, maxAttempt int, backoff FnBackoff) *Backoff {
	if backoff == nil {
		backoff = LinerBackoff(DefaultDurFactor)
	}
	return &Backoff{
		min:        min,
		max:        max,
		maxAttempt: maxAttempt,
		backoff:    backoff,
	}
}

const Stop time.Duration = -1

func (b *Backoff) Next() time.Duration {
	if b.attemptNum >= b.maxAttempt {
		return Stop
	}
	b.attemptNum++
	return b.backoff(b.attemptNum, b.min, b.max)
}

func (b *Backoff) Reset() {
	b.attemptNum = 0
}

func LinerBackoff(factor time.Duration) FnBackoff {
	return func(attemptNum int, min, max time.Duration) time.Duration {
		delay := factor*time.Duration(attemptNum) - 1
		if delay < min {
			delay = min
		}
		if delay > max {
			delay = max
		}
		return delay
	}
}
