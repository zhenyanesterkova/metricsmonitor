package retry

import (
	"context"
	"time"
)

type Action int

const (
	Succeed Action = iota
	Fail
	Retry
)

type RetryPolicy func(err error) Action

type Retrier struct {
	backoff     *Backoff
	retryPolicy RetryPolicy
}

func NewRetrier(backoff *Backoff, retryPolicy RetryPolicy) Retrier {
	if retryPolicy == nil {
		retryPolicy = DefaultRetryPolicy
	}

	return Retrier{
		backoff:     backoff,
		retryPolicy: retryPolicy,
	}
}

func (r Retrier) Run(ctx context.Context, work func() error) error {
	defer r.backoff.Reset()
	for {
		err := work()

		switch r.retryPolicy(err) {
		case Succeed, Fail:
			return err
		case Retry:
			var delay time.Duration
			if delay = r.backoff.Next(); delay == Stop {
				return err
			}
			time.Sleep(delay)
		}
	}
}

func DefaultRetryPolicy(err error) Action {
	if err == nil {
		return Succeed
	}
	return Retry
}
