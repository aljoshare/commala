package validator

import (
	"testing"
)

func TestIsBranchNameWhitelisted_Match(t *testing.T) {
	patterns := []string{"^release-please--.*", "^dependabot/.*"}
	branch := "release-please--branches--main--components--runo"

	matched, pattern, err := IsBranchNameWhitelisted(branch, patterns)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !matched {
		t.Errorf("expected branch %q to match patterns, but did not match", branch)
	}
	if pattern != "^release-please--.*" {
		t.Errorf("expected matched pattern '^release-please--.*', got %q", pattern)
	}
}

func TestIsBranchNameWhitelisted_NoMatch(t *testing.T) {
	patterns := []string{"^release-please--.*"}
	branch := "feature/my-feature"

	matched, pattern, err := IsBranchNameWhitelisted(branch, patterns)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if matched {
		t.Errorf("expected branch %q to not match patterns, but matched %q", branch, pattern)
	}
	if pattern != "" {
		t.Errorf("expected empty pattern, got %q", pattern)
	}

	// Also test empty patterns slice
	matched, _, err = IsBranchNameWhitelisted(branch, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if matched {
		t.Errorf("expected false for empty patterns, got true")
	}
}

func TestIsBranchNameWhitelisted_InvalidRegex(t *testing.T) {
	patterns := []string{"[invalid("}
	branch := "feature/my-feature"

	matched, _, err := IsBranchNameWhitelisted(branch, patterns)
	if err == nil {
		t.Fatalf("expected error for invalid regex, got nil")
	}
	if matched {
		t.Errorf("expected matched to be false on error, got true")
	}
}

func TestIsCommitMessageWhitelisted_Match(t *testing.T) {
	patterns := []string{`^chore\(.*\): release .*`}
	msg := "chore(main): release 0.6.0\n\nRelease notes here"

	matched, pattern, err := IsCommitMessageWhitelisted(msg, patterns)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !matched {
		t.Errorf("expected commit message to match pattern, but did not match")
	}
	if pattern != `^chore\(.*\): release .*` {
		t.Errorf("expected pattern %q, got %q", `^chore\(.*\): release .*`, pattern)
	}
}

func TestIsCommitMessageWhitelisted_NoMatch(t *testing.T) {
	patterns := []string{`^chore\(.*\): release .*`}
	msg := "feat: regular feature without release"

	matched, pattern, err := IsCommitMessageWhitelisted(msg, patterns)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if matched {
		t.Errorf("expected commit message to not match, but matched %q", pattern)
	}

	// Empty patterns
	matched, _, err = IsCommitMessageWhitelisted(msg, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if matched {
		t.Errorf("expected false for empty patterns")
	}
}

func TestIsCommitMessageWhitelisted_InvalidRegex(t *testing.T) {
	patterns := []string{`[invalid(`}
	msg := "chore: foo"

	matched, _, err := IsCommitMessageWhitelisted(msg, patterns)
	if err == nil {
		t.Fatalf("expected error for invalid regex, got nil")
	}
	if matched {
		t.Errorf("expected matched to be false on error")
	}
}
