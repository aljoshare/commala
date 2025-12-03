package validator

import "time"

type ResultMessage struct {
	Valid   bool
	Message string
}

type ValidationResult struct {
	Validator  string
	Valid      bool
	Messages   map[string]ResultMessage
	Summary    string
	Assertions int
	Failures   int
	Duration   time.Duration
}
