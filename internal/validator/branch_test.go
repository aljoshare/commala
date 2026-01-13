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
	result, err := b.Validate(m, []string{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.Valid != true {
		t.Errorf("Expected branch to be conventional, got %v", result.Valid)
	}
}

// Whitelist tests
func TestIsBranchAuthorWhitelisted_EmptyWhitelist(t *testing.T) {
	mockGit := git.MockGit{}

	whitelisted, email, err := IsBranchAuthorWhitelisted([]string{}, mockGit)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if whitelisted {
		t.Error("Expected not whitelisted with empty whitelist")
	}
	if email != "" {
		t.Errorf("Expected empty email, got %s", email)
	}
}

func TestIsBranchAuthorWhitelisted_NoMatch(t *testing.T) {
	mockGit := git.MockGit{}
	whitelist := []string{"dependabot[bot]@users.noreply.github.com"}

	whitelisted, email, err := IsBranchAuthorWhitelisted(whitelist, mockGit)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if whitelisted {
		t.Error("Expected not whitelisted when email doesn't match")
	}
	if email != "" {
		t.Errorf("Expected empty email, got %s", email)
	}
}
