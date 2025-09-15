package validator

import (
	"testing"

	"github.com/aljoshare/commala/internal/git"
)

func TestHasAuthorName(t *testing.T) {
	validator := AuthorNameValidator{
		authorNames: map[string]string{
			"commit1": "John Doe",
			"commit2": "Jane Smith",
		},
	}

	result, err := validator.hasAuthorName()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result["commit1"] {
		t.Errorf("Expected commit1 to have author name")
	}
	if !result["commit2"] {
		t.Errorf("Expected commit2 to have author name")
	}
}

func TestHasNoAuthorName(t *testing.T) {
	validator := AuthorNameValidator{
		authorNames: map[string]string{
			"commit1": "",
		},
	}

	result, err := validator.hasAuthorName()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result["commit1"] {
		t.Errorf("Expected commit3 to not have author name")
	}
}

func TestAuthorNameValidate(t *testing.T) {
	a := AuthorNameValidator{}
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
		t.Errorf("Expected that commits has author name, got %v", result.Valid)
	}
}
