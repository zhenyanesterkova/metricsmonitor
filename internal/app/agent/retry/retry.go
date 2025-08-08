package retry

import (
	"fmt"
	"time"
)

type RequestFunc func() error

func RetryRequest(
	attemptFunc RequestFunc,
	attemptIntervals []string,
) error {
	err := attemptFunc()
	if err == nil {
		return nil
	}

	reqSuccess := false
	for _, interval := range attemptIntervals {
		dur, errParse := time.ParseDuration(interval)
		if errParse != nil {
			return fmt.Errorf("failed to parse interval: %w", errParse)
		}

		time.Sleep(dur)
		err = attemptFunc()
		if err == nil {
			reqSuccess = true
			break
		}
	}

	if !reqSuccess {
		return fmt.Errorf("all retry attempts failed: %w", err)
	}
	return nil
}
