package validator

import (
	"testing"

	"github.com/aljoshare/commala/internal/git"
	"github.com/aljoshare/commala/internal/utils"
)

func TestMessageIsConventional(t *testing.T) {
	m := MessageValidator{}
	m.messages = make(map[string]string)
	m.messages["commit3"] = "feat(scope): add new feature"
	valid, err := m.isConventional()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !utils.AllTrue(valid) {
		t.Errorf("Expected messages to be conventional, got %v", valid)
	}
}

func TestMessageIsNotConventional(t *testing.T) {
	m := MessageValidator{}
	m.messages = make(map[string]string)
	m.messages["commit4"] = "test"
	m.messages["commit5"] = "not: valid"
	m.messages["commit6"] = "feat(): scope not valid"
	valid, err := m.isConventional()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !utils.AllFalse(valid) {
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
