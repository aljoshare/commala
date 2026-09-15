package validator

import (
	"fmt"
	"sort"

	"github.com/aljoshare/commala/internal/config"
	"github.com/aljoshare/commala/internal/git"
)

func Validate(cr *git.CommitRange, g git.Git, c config.Config) ([]*ValidationResult, error) {
	if len(c.WhitelistBranches) > 0 {
		branchName, err := g.GetBranchName()
		if err != nil {
			return nil, err
		}
		matched, pattern, err := IsBranchNameWhitelisted(branchName, c.WhitelistBranches)
		if err != nil {
			return nil, err
		}
		if matched {
			return createSkippedResults(cr, g, c, pattern)
		}
	}

	var results []*ValidationResult

	validators := []struct {
		enabled  bool
		validate func() (*ValidationResult, error)
	}{
		{c.BranchEnabled, func() (*ValidationResult, error) {
			return BranchValidator{}.Validate(g, c.BranchWhitelist, c.BranchPatterns)
		}},
		{c.MessageEnabled, func() (*ValidationResult, error) {
			return MessageValidator{}.Validate(cr, g, c.MessageWhitelist, c.MessageBranchPatterns, c.MessagePatterns)
		}},
		{c.SignOffEnabled, func() (*ValidationResult, error) {
			return SignOffValidator{}.Validate(cr, g, c.SignOffWhitelist, c.SignOffBranchPatterns, c.SignOffPatterns)
		}},
		{c.AuthorNameEnabled, func() (*ValidationResult, error) {
			return AuthorNameValidator{}.Validate(cr, g, c.AuthorNameWhitelist)
		}},
		{c.AuthorEmailEnabled, func() (*ValidationResult, error) {
			return AuthorEmailValidator{}.Validate(cr, g, c.AuthorEmailWhitelist)
		}},
	}

	type result struct {
		vr  *ValidationResult
		err error
	}
	ch := make(chan result, len(validators))
	var tasks int

	for _, v := range validators {
		if v.enabled {
			tasks++
			go func(validate func() (*ValidationResult, error)) {
				vr, err := validate()
				ch <- result{vr, err}
			}(v.validate)
		}
	}

	for i := 0; i < tasks; i++ {
		r := <-ch
		if r.err != nil {
			return nil, r.err
		}
		results = append(results, r.vr)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Validator < results[j].Validator
	})
	return results, nil
}

func createSkippedResults(cr *git.CommitRange, g git.Git, c config.Config, pattern string) ([]*ValidationResult, error) {
	var results []*ValidationResult
	skipMsg := fmt.Sprintf("Skipped (branch whitelisted: %s)", pattern)
	skipReason := fmt.Sprintf("Branch whitelisted: %s", pattern)

	var commits map[string]string
	if cr != nil {
		commits, _ = g.GetCommitMessages(cr.From, cr.To)
	}

	if c.AuthorEmailEnabled {
		vr := &ValidationResult{
			Validator: "Author Email",
			Valid:     true,
			Messages:  make(map[string]ResultMessage),
		}
		if len(commits) > 0 {
			for h := range commits {
				vr.Messages[h] = NewSkippedResultMessageWithReason(skipMsg, skipReason)
			}
			vr.Skipped = len(commits)
			vr.Summary = fmt.Sprintf("All author emails are present (%d skipped)\n", len(commits))
		} else {
			vr.Messages["author_email"] = NewSkippedResultMessageWithReason(skipMsg, skipReason)
			vr.Skipped = 1
			vr.Summary = "Author Email validation skipped (1 skipped)"
		}
		results = append(results, vr)
	}

	if c.AuthorNameEnabled {
		vr := &ValidationResult{
			Validator: "Author Name",
			Valid:     true,
			Messages:  make(map[string]ResultMessage),
		}
		if len(commits) > 0 {
			for h := range commits {
				vr.Messages[h] = NewSkippedResultMessageWithReason(skipMsg, skipReason)
			}
			vr.Skipped = len(commits)
			vr.Summary = fmt.Sprintf("All author names are present (%d skipped)\n", len(commits))
		} else {
			vr.Messages["author_name"] = NewSkippedResultMessageWithReason(skipMsg, skipReason)
			vr.Skipped = 1
			vr.Summary = "Author Name validation skipped (1 skipped)"
		}
		results = append(results, vr)
	}

	if c.BranchEnabled {
		vr := &ValidationResult{
			Validator: "Branch",
			Valid:     true,
			Skipped:   1,
			Summary:   "Branch validation skipped (1 skipped)",
			Messages: map[string]ResultMessage{
				"branch": NewSkippedResultMessageWithReason(skipMsg, skipReason),
			},
		}
		results = append(results, vr)
	}

	if c.MessageEnabled {
		vr := &ValidationResult{
			Validator: "Conventional Message",
			Valid:     true,
			Messages:  make(map[string]ResultMessage),
		}
		if len(commits) > 0 {
			for h := range commits {
				vr.Messages[h] = NewSkippedResultMessageWithReason(skipMsg, skipReason)
			}
			vr.Skipped = len(commits)
			vr.Summary = fmt.Sprintf("All messages are conventional (%d skipped)\n", len(commits))
		} else {
			vr.Messages["message"] = NewSkippedResultMessageWithReason(skipMsg, skipReason)
			vr.Skipped = 1
			vr.Summary = "Message validation skipped (1 skipped)"
		}
		results = append(results, vr)
	}

	if c.SignOffEnabled {
		vr := &ValidationResult{
			Validator: "Signed-Off Message",
			Valid:     true,
			Messages:  make(map[string]ResultMessage),
		}
		if len(commits) > 0 {
			for h := range commits {
				vr.Messages[h] = NewSkippedResultMessageWithReason(skipMsg, skipReason)
			}
			vr.Skipped = len(commits)
			vr.Summary = fmt.Sprintf("All commits are signed off (%d skipped)\n", len(commits))
		} else {
			vr.Messages["signoff"] = NewSkippedResultMessageWithReason(skipMsg, skipReason)
			vr.Skipped = 1
			vr.Summary = "Sign-off validation skipped (1 skipped)"
		}
		results = append(results, vr)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Validator < results[j].Validator
	})
	return results, nil
}
