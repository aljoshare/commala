package validator

import (
	"testing"

	"github.com/aljoshare/commala/internal/git"
	"github.com/aljoshare/commala/internal/utils"
)

func TestMessageIsSignedOff(t *testing.T) {
	m := SignOffValidator{}
	m.messages = append(m.messages, "feat(scope): add new feature \nSigned-off-by: Johen Doe <john@doe.com>")
	valid, err := m.isSignedOff()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !utils.AllTrueSlice(valid) {
		t.Errorf("Expected messages to be signed off, got %v", valid)
	}
}

func TestMessageIsNotSignedOff(t *testing.T) {
	m := SignOffValidator{}
	m.messages = append(m.messages, "test")
	m.messages = append(m.messages, "not: valid")
	m.messages = append(m.messages, "feat(): scope not valid")
	valid, err := m.isSignedOff()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !utils.AllFalseSlice(valid) {
		t.Errorf("Expected messages to be not signed off, got %v", valid)
	}
}

func TestSignOffValidate(t *testing.T) {
	v := SignOffValidator{}
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
		t.Errorf("Expected message to be signed off, got %v", result.Valid)
	}
}
