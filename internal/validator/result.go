package validator

import "time"

type ValidationResult struct {
	Validator string
	Valid     bool
	Message   []string
	Duration  time.Duration
}
