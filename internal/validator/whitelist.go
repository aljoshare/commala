package validator

import (
	"fmt"
	"regexp"

	"github.com/aljoshare/commala/internal/git"
)

// IsWhitelisted checks if the commit author email is in the whitelist.
// Returns true if whitelisted, along with the matched email.
func IsWhitelisted(commitHash string, whitelist []string, g git.Git, cr *git.CommitRange) (bool, string, error) {
	if len(whitelist) == 0 {
		return false, "", nil
	}

	// Get author emails for the commit range
	authorEmails, err := g.GetCommitAuthorEmails(cr.From, cr.To)
	if err != nil {
		return false, "", err
	}

	// Get the email for this specific commit
	authorEmail, exists := authorEmails[commitHash]
	if !exists {
		return false, "", nil
	}

	// Check if email is in whitelist (exact match)
	for _, whitelistedEmail := range whitelist {
		if authorEmail == whitelistedEmail {
			return true, authorEmail, nil
		}
	}

	return false, "", nil
}

// IsBranchAuthorWhitelisted checks if the current branch's HEAD commit author is whitelisted.
// Special function for BranchValidator since it doesn't process individual commits.
func IsBranchAuthorWhitelisted(whitelist []string, g git.Git) (bool, string, error) {
	if len(whitelist) == 0 {
		return false, "", nil
	}

	// Get HEAD commit
	headCommit, err := g.GetLatestCommit()
	if err != nil {
		return false, "", err
	}

	// Create a single-commit range to get author
	authorEmails, err := g.GetCommitAuthorEmails(headCommit, headCommit)
	if err != nil {
		return false, "", err
	}

	authorEmail, exists := authorEmails[headCommit]
	if !exists {
		return false, "", nil
	}

	// Check if email is in whitelist
	for _, whitelistedEmail := range whitelist {
		if authorEmail == whitelistedEmail {
			return true, authorEmail, nil
		}
	}

	return false, "", nil
}

// IsBranchNameWhitelisted checks if the branch name matches any regex in patterns.
// Returns true if whitelisted, along with the matched pattern.
func IsBranchNameWhitelisted(branchName string, patterns []string) (bool, string, error) {
	if len(patterns) == 0 {
		return false, "", nil
	}

	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return false, "", fmt.Errorf("invalid branch regex pattern %q: %w", pattern, err)
		}
		if re.MatchString(branchName) {
			return true, pattern, nil
		}
	}

	return false, "", nil
}

// IsCommitMessageWhitelisted checks if the commit message matches any regex in patterns.
// Returns true if whitelisted, along with the matched pattern.
func IsCommitMessageWhitelisted(message string, patterns []string) (bool, string, error) {
	if len(patterns) == 0 {
		return false, "", nil
	}

	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return false, "", fmt.Errorf("invalid commit message regex pattern %q: %w", pattern, err)
		}
		if re.MatchString(message) {
			return true, pattern, nil
		}
	}

	return false, "", nil
}
