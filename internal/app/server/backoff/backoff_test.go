package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	minDelay   = time.Second
	maxDelay   = 5 * time.Second
	maxAttempt = 3
)

func TestBackoff(t *testing.T) {
	b := New(
		minDelay,
		maxDelay,
		maxAttempt,
	)

	t.Run("success - Next(): ", func(t *testing.T) {
		wantB := &Backoff{
			min:        minDelay,
			max:        maxDelay,
			maxAttempt: maxAttempt,
			attemptNum: 1,
			nextDelay:  minDelay + 2*time.Second,
		}
		del := b.Next()
		require.Equal(t, minDelay, del)
		require.Equal(t, b, wantB)
	})

	b.Next()
	b.Next()

	t.Run("stop - Next(): ", func(t *testing.T) {
		wantB := &Backoff{
			min:        minDelay,
			max:        maxDelay,
			maxAttempt: maxAttempt,
			attemptNum: 3,
			nextDelay:  7 * time.Second,
		}
		del := b.Next()
		require.Equal(t, Stop, del)
		require.Equal(t, b, wantB)
	})

	t.Run("success - Reset(): ", func(t *testing.T) {
		wantB := &Backoff{
			min:        minDelay,
			max:        maxDelay,
			maxAttempt: maxAttempt,
			attemptNum: 0,
			nextDelay:  minDelay,
		}
		b.Reset()
		require.Equal(t, b, wantB)
	})
}
