package utils

import (
	"time"
)

func Duration(start time.Time) time.Duration {
	elapsed := time.Since(start)
	return elapsed
}
