package validator

import (
	"testing"

	"github.com/aljoshare/commala/internal/git"
)

func TestHasAuthorEmail(t *testing.T) {
	validator := AuthorEmailValidator{
		authorEmails: map[string]string{
			"commit1": "john@doe.com",
			"commit2": "jane@smith.com",
		},
	}

	result, err := validator.hasAuthorEmail()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result["commit1"] {
		t.Errorf("Expected commit1 to have author email")
	}
	if !result["commit2"] {
		t.Errorf("Expected commit2 to have author email")
	}
}

func TestHasNoAuthorEmail(t *testing.T) {
	validator := AuthorEmailValidator{
		authorEmails: map[string]string{
			"commit1": "",
		},
	}

	result, err := validator.hasAuthorEmail()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result["commit1"] {
		t.Errorf("Expected commit3 to not have author email")
	}
}

func TestAuthorEmailValidate(t *testing.T) {
	a := AuthorEmailValidator{}
	m := git.MockGit{}
	r := git.CommitRange{
		From: "commit1",
		To:   "commit2",
	}
	result, err := a.Validate(&r, m)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.Valid != true {
		t.Errorf("Expected that commit has author email, got %v", result.Valid)
	}
}
