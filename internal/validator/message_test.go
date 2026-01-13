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
	result, err := v.Validate(&r, m, []string{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.Valid != true {
		t.Errorf("Expected message to be conventional, got %v", result.Valid)
	}
}

// Whitelist tests
func TestIsWhitelisted_EmptyWhitelist(t *testing.T) {
	mockGit := git.MockGit{}
	cr := &git.CommitRange{From: "commit1", To: "commit1"}

	whitelisted, email, err := IsWhitelisted("commit1", []string{}, mockGit, cr)

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

func TestIsWhitelisted_ExactMatch(t *testing.T) {
	mockGit := git.MockGit{}
	cr := &git.CommitRange{From: "commit1", To: "commit1"}
	whitelist := []string{"test@user.de"}

	whitelisted, email, err := IsWhitelisted("commit1", whitelist, mockGit, cr)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !whitelisted {
		t.Error("Expected whitelisted for exact match")
	}
	if email != "test@user.de" {
		t.Errorf("Expected matched email 'test@user.de', got %s", email)
	}
}

func TestIsWhitelisted_NoMatch(t *testing.T) {
	mockGit := git.MockGit{}
	cr := &git.CommitRange{From: "commit1", To: "commit1"}
	whitelist := []string{"dependabot[bot]@users.noreply.github.com"}

	whitelisted, email, err := IsWhitelisted("commit1", whitelist, mockGit, cr)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if whitelisted {
		t.Error("Expected not whitelisted for no match")
	}
	if email != "" {
		t.Errorf("Expected empty email, got %s", email)
	}
}

func TestIsWhitelisted_MultipleWhitelistEntries(t *testing.T) {
	mockGit := git.MockGit{}
	cr := &git.CommitRange{From: "commit2", To: "commit2"}
	whitelist := []string{
		"dependabot[bot]@users.noreply.github.com",
		"anothertest@user.de",
	}

	whitelisted, email, err := IsWhitelisted("commit2", whitelist, mockGit, cr)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !whitelisted {
		t.Error("Expected whitelisted when matching second entry")
	}
	if email != "anothertest@user.de" {
		t.Errorf("Expected matched email 'anothertest@user.de', got %s", email)
	}
}
