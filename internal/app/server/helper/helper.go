package helper

import (
	"fmt"
	"time"
)

func GetDurationsFromArrString(arr []string) ([]time.Duration, error) {
	durs := make([]time.Duration, 0)
	for _, stringDur := range arr {
		dur, err := time.ParseDuration(stringDur)
		if err != nil {
			return durs, fmt.Errorf("failed parse duration from string: %w", err)
		}
		durs = append(durs, dur)
	}
	return durs, nil
}
