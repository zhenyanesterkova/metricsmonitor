package config

import "time"

const (
	DefaultMinRetryDelay   = time.Second
	DefaultMaxRetryDelay   = 5 * time.Second
	DefaultMaxRetryAttempt = 3
)

type RetryConfig struct {
	MinDelay   time.Duration `json:"min_delay"`
	MaxDelay   time.Duration `json:"max_delay"`
	MaxAttempt int           `json:"max_attempt"`
}
