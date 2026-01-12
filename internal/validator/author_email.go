package validator

import (
	"fmt"
	"time"

	"github.com/aljoshare/commala/internal/git"
	"github.com/aljoshare/commala/internal/utils"
	log "github.com/sirupsen/logrus"
)

type AuthorEmailValidator struct {
	authorEmails map[string]string
}

func (m AuthorEmailValidator) Validate(cr *git.CommitRange, g git.Git, whitelist []string) (*ValidationResult, error) {
	log.Debugf("Validating if commits has author emails from %s to %s", cr.From, cr.To)
	vr := ValidationResult{
		Validator:  "Author Email",
		Assertions: 0,
		Failures:   0,
		Skipped:    0,
	}
	defer func() {
		vr.Duration = utils.Duration(time.Now())
	}()
	var err error
	m.authorEmails, err = g.GetCommitAuthorEmails(cr.From, cr.To)
	if err != nil {
		return nil, err
	}
	vm, err := m.hasAuthorEmail()
	if err != nil {
		return nil, err
	}

	vr.Messages = make(map[string]ResultMessage, len(m.authorEmails))
	for commitHash := range m.authorEmails {
		// Check whitelist before validation
		whitelisted, matchedEmail, err := IsWhitelisted(commitHash, whitelist, g, cr)
		if err != nil {
			return nil, err
		}

		if whitelisted {
			vr.Skipped++
			vr.Messages[commitHash] = NewSkippedResultMessage(matchedEmail)
			continue
		}

		vr.Assertions++
		if !vm[commitHash] {
			vr.Messages[commitHash] = ResultMessage{Valid: false, Message: fmt.Sprintf("Commit Message: \"%s\" has no author email\n", commitHash)}
			vr.Failures++
		} else {
			vr.Messages[commitHash] = ResultMessage{Valid: true, Message: fmt.Sprintf("Commit Message: \"%s\" has author email\n", commitHash)}
		}
	}

	// Update summary with skip count
	if vr.Failures == 0 {
		vr.Valid = true
		if vr.Skipped > 0 {
			vr.Summary = fmt.Sprintf("All author emails are present (%d skipped)\n", vr.Skipped)
		} else {
			vr.Summary = "All commit messages have author email set\n"
		}
	} else {
		vr.Valid = true
		totalChecked := vr.Assertions + vr.Skipped
		passedCount := vr.Assertions - vr.Failures
		if vr.Skipped > 0 {
			vr.Summary = fmt.Sprintf("%d of %d author emails are present (%d skipped)\n", passedCount, totalChecked, vr.Skipped)
		} else {
			vr.Summary = "Not all commit messages have author email set\n"
		}
	}
	return &vr, nil
}

func (m AuthorEmailValidator) hasAuthorEmail() (map[string]bool, error) {
	am := make(map[string]bool, len(m.authorEmails))
	for ch, an := range m.authorEmails {
		if an != "" {
			am[ch] = true
		} else {
			am[ch] = false
		}
	}
	return am, nil
}
