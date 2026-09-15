package validator

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aljoshare/commala/internal/git"
	"github.com/aljoshare/commala/internal/utils"
	log "github.com/sirupsen/logrus"
)

type MessageValidator struct {
	messages map[string]string
}

func (m MessageValidator) Validate(cr *git.CommitRange, g git.Git, emailWhitelist []string, branchPatterns []string, messagePatterns []string) (*ValidationResult, error) {
	log.Debugf("Validating if commit messages are conventional from %s to %s", cr.From, cr.To)
	vr := ValidationResult{
		Validator: "Conventional Message",
		Skipped:   0,
	}
	defer func() {
		vr.Duration = utils.Duration(time.Now())
	}()
	var err error
	m.messages, err = g.GetCommitMessages(cr.From, cr.To)
	if err != nil {
		return nil, err
	}

	vr.Messages = make(map[string]ResultMessage, len(m.messages))

	// 1. Check if current branch matches branchPatterns
	if len(branchPatterns) > 0 {
		branchName, err := g.GetBranchName()
		if err != nil {
			return nil, err
		}
		branchWhitelisted, matchedPattern, err := IsBranchNameWhitelisted(branchName, branchPatterns)
		if err != nil {
			return nil, err
		}
		if branchWhitelisted {
			vr.Valid = true
			for commitHash := range m.messages {
				vr.Skipped++
				vr.Messages[commitHash] = NewSkippedResultMessageWithReason(
					fmt.Sprintf("Skipped (branch pattern matched: %s)", matchedPattern),
					fmt.Sprintf("Branch pattern matched: %s", matchedPattern),
				)
			}
			if vr.Skipped > 0 {
				vr.Summary = fmt.Sprintf("All messages are conventional (%d skipped)\n", vr.Skipped)
			} else {
				vr.Summary = "All commit messages are conventional\n"
			}
			return &vr, nil
		}
	}

	// 2. Process each commit
	for commitHash, commitMsg := range m.messages {
		// Check email whitelist before validation
		whitelisted, authorEmail, err := IsWhitelisted(commitHash, emailWhitelist, g, cr)
		if err != nil {
			return nil, err
		}

		if whitelisted {
			vr.Skipped++
			vr.Messages[commitHash] = NewSkippedResultMessage(authorEmail)
			continue
		}

		// Check commit message against messagePatterns
		msgWhitelisted, matchedPattern, err := IsCommitMessageWhitelisted(commitMsg, messagePatterns)
		if err != nil {
			return nil, err
		}

		if msgWhitelisted {
			vr.Skipped++
			vr.Messages[commitHash] = NewSkippedResultMessageWithReason(
				fmt.Sprintf("Skipped (commit message pattern matched: %s)", matchedPattern),
				fmt.Sprintf("Commit message pattern matched: %s", matchedPattern),
			)
			continue
		}

		// Existing validation logic
		vr.Assertions++
		matched, err := regexp.Match(`^(build|chore|ci|docs|feat|fix|perf|refactor|revert|style|test){1}(\([\w\-\.]+\))?(!)?: ([\w ])+([\s\S]*)`, []byte(commitMsg))
		if err != nil {
			return nil, err
		}
		if !matched {
			vr.Messages[commitHash] = ResultMessage{Valid: false, Message: fmt.Sprintf("Commit Message: \"%s\" is not conventional\n", commitHash)}
			vr.Failures++
		} else {
			vr.Messages[commitHash] = ResultMessage{Valid: true, Message: fmt.Sprintf("Commit Message: \"%s\" is conventional\n", commitHash)}
		}
	}

	// Update summary with skip count
	if vr.Failures == 0 {
		vr.Valid = true
		if vr.Skipped > 0 {
			vr.Summary = fmt.Sprintf("All messages are conventional (%d skipped)\n", vr.Skipped)
		} else {
			vr.Summary = "All commit messages are conventional\n"
		}
	} else {
		vr.Valid = false
		totalChecked := vr.Assertions + vr.Skipped
		passedCount := vr.Assertions - vr.Failures
		if vr.Skipped > 0 {
			vr.Summary = fmt.Sprintf("%d of %d messages are conventional (%d skipped)\n", passedCount, totalChecked, vr.Skipped)
		} else {
			vr.Summary = "Not all commit messages are conventional\n"
		}
	}
	return &vr, nil
}

func (m MessageValidator) isConventional() (map[string]bool, error) {
	var l = make(map[string]bool, len(m.messages))
	re := regexp.MustCompile(`^(build|chore|ci|docs|feat|fix|perf|refactor|revert|style|test){1}(\([\w\-\.]+\))?(!)?: ([\w ])+([\s\S]*)`)
	for i, cm := range m.messages {
		matched := re.Match([]byte(cm))
		l[i] = matched
	}
	return l, nil
}
