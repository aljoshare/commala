package validator

import (
	"testing"

	"github.com/aljoshare/commala/internal/config"
	"github.com/aljoshare/commala/internal/git"
)

func TestValidate_GlobalBranchWhitelist(t *testing.T) {
	m := git.MockGit{
		BranchName: "release-please--branches--main--components--runo",
		CommitMessages: map[string]string{
			"c1": "chore(main): release 0.6.0",
		},
	}
	cr := &git.CommitRange{From: "c1", To: "c1"}

	// Configure with global branch whitelist matching the branch name
	cfg := config.Config{
		BranchEnabled:      true,
		MessageEnabled:     true,
		SignOffEnabled:     true,
		AuthorNameEnabled:  true,
		AuthorEmailEnabled: true,
		WhitelistBranches:  []string{"^release-please--.*"},
	}

	results, err := Validate(cr, m, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 5 {
		t.Fatalf("expected 5 validation results, got %d", len(results))
	}

	for _, res := range results {
		if !res.Valid {
			t.Errorf("validator %s expected to be valid, but failed", res.Validator)
		}
		if res.Skipped == 0 {
			t.Errorf("validator %s expected to have skipped > 0", res.Validator)
		}
		if res.Failures != 0 {
			t.Errorf("validator %s expected 0 failures, got %d", res.Validator, res.Failures)
		}
		if len(res.Messages) == 0 {
			t.Errorf("validator %s expected to have messages, got 0", res.Validator)
		}
		for _, msg := range res.Messages {
			if !msg.Skipped {
				t.Errorf("validator %s message expected to be skipped: %+v", res.Validator, msg)
			}
			if msg.SkipReason != "Branch whitelisted: ^release-please--.*" {
				t.Errorf("validator %s message expected SkipReason 'Branch whitelisted: ^release-please--.*', got %q", res.Validator, msg.SkipReason)
			}
		}
	}
}

func TestValidate_GlobalBranchWhitelist_NoMatch(t *testing.T) {
	m := git.MockGit{
		BranchName: "invalid-branch-name",
		CommitMessages: map[string]string{
			"c1": "not conventional",
		},
	}
	cr := &git.CommitRange{From: "c1", To: "c1"}

	cfg := config.Config{
		BranchEnabled:     true,
		MessageEnabled:    true,
		WhitelistBranches: []string{"^release-please--.*"},
	}

	results, err := Validate(cr, m, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Because branch doesn't match and commit is not conventional, validators should run and fail
	var branchFailed, messageFailed bool
	for _, res := range results {
		if res.Validator == "Branch" && !res.Valid {
			branchFailed = true
		}
		if res.Validator == "Conventional Message" && !res.Valid {
			messageFailed = true
		}
	}

	if !branchFailed {
		t.Errorf("expected branch validator to fail when global branch pattern did not match")
	}
	if !messageFailed {
		t.Errorf("expected message validator to fail when global branch pattern did not match")
	}
}
