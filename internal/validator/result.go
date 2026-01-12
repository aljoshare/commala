package validator

import (
	"fmt"
	"time"
)

type ResultMessage struct {
	Valid      bool
	Message    string
	Skipped    bool   // true if validation was skipped
	SkipReason string // explanation for skip (e.g., "Author whitelisted: email@example.com")
}

type ValidationResult struct {
	Validator  string
	Valid      bool
	Messages   map[string]ResultMessage
	Summary    string
	Assertions int
	Failures   int
	Skipped    int // count of skipped validations
	Duration   time.Duration
}

// NewSkippedResultMessage creates a ResultMessage for a skipped validation
func NewSkippedResultMessage(whitelistedEmail string) ResultMessage {
	return ResultMessage{
		Valid:      true, // Skipped commits don't fail validation
		Skipped:    true,
		Message:    fmt.Sprintf("Skipped (author whitelisted: %s)", whitelistedEmail),
		SkipReason: fmt.Sprintf("Author whitelisted: %s", whitelistedEmail),
	}
}
