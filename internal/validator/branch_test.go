package validator

import (
	"testing"

	"github.com/aljoshare/commala/internal/git"
)

func TestBranchIsConventional(t *testing.T) {
	b := BranchValidator{}
	b.branchName = "feature/awesome-feature"
	valid, err := b.isConventional()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !valid {
		t.Errorf("Expected branch to be conventional, got %v", valid)
	}
}

func TestBranchIsNotConventional(t *testing.T) {
	b := BranchValidator{}
	b.branchName = "main"
	valid, err := b.isConventional()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if valid {
		t.Errorf("Expected branch to be not conventional, got %v", valid)
	}
}

func TestBranchValidate(t *testing.T) {
	b := BranchValidator{}
	m := git.MockGit{}
	result, err := b.Validate(m)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.Valid != true {
		t.Errorf("Expected branch to be conventional, got %v", result.Valid)
	}
}
