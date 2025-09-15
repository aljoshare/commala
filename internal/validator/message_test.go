package validator

import (
	"testing"

	"github.com/aljoshare/commala/internal/git"
	"github.com/aljoshare/commala/internal/utils"
)

func TestMessageIsConventional(t *testing.T) {
	m := MessageValidator{}
	m.messages = append(m.messages, "feat(scope): add new feature")
	valid, err := m.isConventional()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !utils.AllTrueSlice(valid) {
		t.Errorf("Expected messages to be conventional, got %v", valid)
	}
}

func TestMessageIsNotConventional(t *testing.T) {
	m := MessageValidator{}
	m.messages = append(m.messages, "test")
	m.messages = append(m.messages, "not: valid")
	m.messages = append(m.messages, "feat(): scope not valid")
	valid, err := m.isConventional()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !utils.AllFalseSlice(valid) {
		t.Errorf("Expected messages to be not conventional, got %v", valid)
	}
}

func TestMessageValidate(t *testing.T) {
	v := MessageValidator{}
	m := git.MockGit{}
	r := git.CommitRange{
		From: "commit1",
		To:   "commit2",
	}
	result, err := v.Validate(&r, m)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.Valid != true {
		t.Errorf("Expected message to be conventional, got %v", result.Valid)
	}
}
