package config

import "time"

const (
	DefaultMinDelay   = time.Second
	DefaultMaxDelay   = 5 * time.Second
	DefaultMaxAttempt = 4
)

type RetryConfig struct {
	Min        time.Duration
	Max        time.Duration
	MaxAttempt int
}
